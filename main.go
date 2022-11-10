package main

import (
	initializer "github.com/chumvan/confdb/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	config, err := initializer.LoadEnv(".")
	if err != nil {
		panic(err)
	}
	_, err = initializer.ConnectToDB(config.DBDns)
	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
	r.Run()
}
