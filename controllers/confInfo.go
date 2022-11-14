package controller

import (
	"fmt"
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
	fmt.Printf("Users type: %#v\n Users: %#v", confInfo.Users, confInfo.Users)
	err = model.AddNewConfInfo(&confInfo)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": confInfo})
}
