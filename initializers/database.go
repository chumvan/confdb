package initializer

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB(dns string) (db *gorm.DB, err error) {
	DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	fmt.Println("Connected to DB")

	return DB, nil
}
