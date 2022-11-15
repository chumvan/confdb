package router

import (
	controller "github.com/chumvan/confdb/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() (r *gin.Engine) {
	r = gin.Default()

	r.GET("/confInfos", controller.GetConfInfos)
	r.POST("/confInfos", controller.CreateAConfInfo)

	return r
}
