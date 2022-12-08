package initializer

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() (db *gorm.DB, err error) {
	dsn := "host=10.5.0.9 user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Helsinki"
	fmt.Printf("dsn: %s\n", dsn)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	fmt.Println("Connected to DB")

	return DB, nil
}
