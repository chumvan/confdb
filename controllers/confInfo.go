package controller

import (
	"net/http"

	model "github.com/chumvan/confdb/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ConfInfoInput struct {
	ConfUri datatypes.URL
	Subject string
	Users   []model.User
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
		c.JSON(http.StatusBadRequest, gin.H{"mesage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Data": confInfo})
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
		Users:   input.Users,
	}

	err = model.AddNewConfInfo(&confInfo)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Data": confInfo})
}

type InputUser struct {
	EntityUrl datatypes.URL
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

func DeleteUserFromTopic(c *gin.Context) {
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

		if err := model.DeleteUserFromTopic(topic, user, &confInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "topic not found"})
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
