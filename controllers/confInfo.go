package controller

import (
	"net"
	"net/http"
	"os"
	"strconv"

	model "github.com/chumvan/confdb/models"
	"github.com/gin-gonic/gin"
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
	ConfInfo  model.ConfInfo
	TopicIP   net.IP
	TopicPort int
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
	forwarderRtpInPort, _ := strconv.Atoi(forwarderPortStr)

	resp := ResponseCreateConfInfo{
		ConfInfo:  confInfo,
		TopicIP:   net.IP(forwarderIP),
		TopicPort: forwarderRtpInPort,
	}

	c.JSON(http.StatusCreated, gin.H{"Data": resp})
}

type InputUser struct {
	EntityUrl string
	Role      string
}

func AddUserToConfInfo(c *gin.Context) {
	topic := c.Params.ByName("topic")
	if topic != "" {
		var confInfo model.ConfInfo
		var inputUser InputUser
		if err := c.ShouldBindJSON(&inputUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		user := model.User{
			EntityUrl: inputUser.EntityUrl,
			Role:      inputUser.Role,
		}

		if err := model.PatchUserToTopic(topic, user, &confInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		users := confInfo.Users
		c.JSON(http.StatusOK, users)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "topic not found"})

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
