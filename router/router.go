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

		topicMode := v1.Group("/topicMode")
		{
			topicMode.GET("/confInfos/:topic", controller.GetTopicInfo)
			topicMode.PATCH("/confInfos/:topic", controller.AddUserToConfInfo)
			topicMode.DELETE("/confInfos/:topic", controller.DeleteUserFromTopic)
		}
	}
	return r
}
