package main

import (
	"log"

	initializer "github.com/chumvan/confdb/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	_, err := initializer.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
	r.Run()
}
