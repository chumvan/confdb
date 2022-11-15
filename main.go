package main

import (
	"log"

	initializer "github.com/chumvan/confdb/initializers"
	model "github.com/chumvan/confdb/models"
	"github.com/chumvan/confdb/router"
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

	r := router.SetupRouter()

	if err = model.AddMockData(); err != nil {
		log.Fatal(err)
	}

	r.Run()
}
