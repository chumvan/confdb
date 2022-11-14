package main

import (
	"log"

	controller "github.com/chumvan/confdb/controllers"
	initializer "github.com/chumvan/confdb/initializers"
	model "github.com/chumvan/confdb/models"
	"github.com/gin-gonic/gin"
)

func main() {
	DB, err := initializer.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	err = DB.AutoMigrate(&model.ConfInfo{}, model.User{})
	if err != nil {
		log.Fatal("failed to migrate", err)
	}

	r := gin.Default()

	r.GET("/confInfos", controller.GetConfInfos)
	r.POST("/confInfos", controller.CreateAConfInfo)
	r.GET("/confInfos/:id", controller.GetAConfInfoById)

	r.Run()
}
