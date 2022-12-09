package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	model "github.com/chumvan/confdb/models"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

type ConfInfoInput struct {
	ConfUri string
	Subject string
	Creator model.User
}

func GetConfInfos(c *gin.Context) {
	var confInfos []model.ConfInfo
	err := model.GetAllConfInfos(&confInfos)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// Filter by topic
	topic, ok := c.GetQuery("topic")
	if ok {
		var confInfo model.ConfInfo
		if err := model.GetConfInfoByTopic(topic, &confInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"mesage": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Data": confInfo})
		return
	}

	// Default return all ConfInfo
	c.JSON(http.StatusAccepted, gin.H{"Data": confInfos})
}

func GetConfInfoById(c *gin.Context) {
	// Filter by Id
	id := c.Params.ByName("id")

	var confInfo model.ConfInfo
	if err := model.GetConfInfoById(id, &confInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Data": confInfo})
}

type ResponseCreateConfInfo struct {
	ConfUri   string
	TopicIP   string
	TopicPort string
}

func CreateAConfInfo(c *gin.Context) {
	var input ConfInfoInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	confInfo := model.ConfInfo{
		ConfUri: input.ConfUri,
		Subject: input.Subject,
		Users: []model.User{
			input.Creator,
		},
	}

	err = model.AddNewConfInfo(&confInfo)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	forwarderIP := os.Getenv("FORWARDER_IP")
	forwarderPortStr := os.Getenv("FORWARDER_RTP_IN_PORT")
	fmt.Printf("forwarderIP: %s\nforwarderPort: %s\n", forwarderIP, forwarderPortStr)

	resp := ResponseCreateConfInfo{
		ConfUri:   confInfo.ConfUri,
		TopicIP:   forwarderIP,
		TopicPort: forwarderPortStr,
	}
	fmt.Printf("response from confDB: %v\n", resp)
	c.JSON(http.StatusCreated, gin.H{"Data": resp})
}

type InputUser struct {
	EntityUrl string
	Role      string
	PortRTP   int
}

func AddUserToConfInfo(c *gin.Context) {
	topic := c.Params.ByName("topic")
	if topic != "" {
		var found model.ConfInfo
		var inputUser InputUser
		if err := c.ShouldBindJSON(&inputUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// check if topic has been created
		err := model.GetConfInfoByTopic(topic, &found)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "topic not found"})
			return
		}
		fmt.Printf("users in found topic: %v\n", found.Users)

		user := model.User{
			EntityUrl: inputUser.EntityUrl,
			Role:      inputUser.Role,
			PortRTP:   inputUser.PortRTP,
		}

		if err := model.PatchUserToTopic(topic, user, &found); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		users := found.Users

		fmt.Printf("users after patch: %v\n", users)

		// update in topic (forwarder)
		forwarderIP := os.Getenv("FORWARDER_IP")
		forwarderRESTPortStr := os.Getenv("FORWARDER_REST_PORT")

		topicUrl := fmt.Sprintf("http://%s:%s/users", forwarderIP, forwarderRESTPortStr)
		fmt.Printf("topicUrl: %s\n", topicUrl)
		payload, err := json.Marshal(users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to parse user-slice"})
			return
		}
		fmt.Printf("users: %v", users)
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, topicUrl, bytes.NewBuffer(payload))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "fail to create PUT req"})
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "fail to send PUT req"})
			return
		}

		fmt.Printf("updated user at topic: %v", *resp)

		defer req.Body.Close()
		// end of update in topic
		// reply to client
		c.JSON(http.StatusOK, users)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "empty topic string"})

}

// func GetUserFromTopic(c *gin.Context) {
// }

func DeleteUserFromTopic(c *gin.Context) {
	topic, ok := c.Params.Get("topic")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "empty topic"})
		return
	}
	userId, ok := c.Params.Get("userId")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "empty userId"})
		return
	}

	if err := model.DeleteUserFromTopic(topic, userId, &model.ConfInfo{}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// update in topic (forwarder)
	found := &model.ConfInfo{}
	err := model.GetConfInfoByTopic(topic, found)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "topic not found"})
	}
	users := found.Users
	forwarderIP := os.Getenv("FORWARDER_IP")
	forwarderRESTPortStr := os.Getenv("FORWARDER_REST_PORT")

	topicUrl := fmt.Sprintf("http://%s:%s/users", forwarderIP, forwarderRESTPortStr)
	fmt.Printf("topicUrl: %s\n", topicUrl)
	payload, err := json.Marshal(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to parse user-slice"})
		return
	}
	fmt.Printf("users: %v", users)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, topicUrl, bytes.NewBuffer(payload))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "fail to create PUT req"})
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "fail to send PUT req"})
		return
	}

	fmt.Printf("updated user at topic: %v", *resp)

	defer req.Body.Close()
	// end of update in topic
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func GetTopicInfo(c *gin.Context) {
	topic := c.Params.ByName("topic")
	if topic != "" {
		var confInfo model.ConfInfo
		if err := model.GetTopicInfo(topic, &confInfo); err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"Data": confInfo})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "topic not found"})
}
