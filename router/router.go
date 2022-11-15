package router

import (
	controller "github.com/chumvan/confdb/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() (r *gin.Engine) {
	r = gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/confInfos/:id", controller.GetConfInfoById)
		v1.GET("/confInfos", controller.GetConfInfos)
		v1.POST("/confInfos", controller.CreateAConfInfo)

		v1.GET("/topicMode/confInfos/:topic", controller.GetTopicInfo)
		v1.PATCH("/topicMode/confInfos/:topic", controller.AddUserToConfInfo)
		v1.DELETE("/topicMode/confInfos/:topic", controller.DeleteUserFromTopic)
	}

	return r
}
