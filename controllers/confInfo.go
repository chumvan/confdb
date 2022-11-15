package controller

import (
	"net/http"

	model "github.com/chumvan/confdb/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func GetConfInfos(c *gin.Context) {
	var confInfos []model.ConfInfo
	err := model.GetAllConfInfos(&confInfos)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// Filter by Id
	id, ok := c.GetQuery("id")
	if ok {
		var confInfo model.ConfInfo
		if err := model.GetOneConfInfoById(id, &confInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"mesage": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": confInfo})
		return
	}

	// Filter by topic
	topic, ok := c.GetQuery("topic")
	if ok {
		var confInfo model.ConfInfo
		if err := model.GetOneConfInfoByTopic(topic, &confInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"mesage": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": confInfo})
		return
	}

	// Default return all ConfInfo
	c.JSON(http.StatusAccepted, gin.H{"data": confInfos})
}

type ConfInfoInput struct {
	ConfUri datatypes.URL
	Subject string
	Users   []model.User
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
	c.JSON(http.StatusCreated, gin.H{"data": confInfo})
}
