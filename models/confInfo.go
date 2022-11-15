package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	initializer "github.com/chumvan/confdb/initializers"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	UserID    uint `gorm:"primary_key"`
	EntityUrl datatypes.URL
	Role      string
	ConfRefer string
}

type ConfInfo struct {
	gorm.Model
	ConfUri datatypes.URL
	Subject string
	Users   []User `gorm:"foreignKey:ConfRefer"`
}

// GET /confInfos
// Get all ConfInfos
func GetAllConfInfos(confInfos *[]ConfInfo) (err error) {
	if err = initializer.DB.Preload("Users").Find(&confInfos).Error; err != nil {
		return err
	}
	return nil
}

// GET /confInfos/:id
// Get a ConfInfo by its id
func GetOneConfInfoById(confId string, confInfo *ConfInfo) (err error) {
	if err = initializer.DB.Where("id = ?", confId).First(&confInfo).Error; err != nil {
		return err
	}
	return nil
}

// POST /confInfos
// Create a ConfInfo in
func AddNewConfInfo(confInfo *ConfInfo) (err error) {
	if err = initializer.DB.Create(&confInfo).Error; err != nil {
		return err
	}
	return nil
}

type InputConfInfos struct {
	Data []ConfInfo
}

type InputConfInfo struct {
	ConfUri datatypes.URL
	Subject string
	Users   []User
}

// Add initialized mocked data for development
func AddMockData() (err error) {
	jsonFile, err := os.Open("./initializers/data.json")
	if err != nil {
		return err
	}
	fmt.Println("json file read")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var inputConfInfos InputConfInfos
	err = json.Unmarshal(byteValue, &inputConfInfos)
	if err != nil {
		return err
	}

	for _, in := range inputConfInfos.Data {
		confInfo := ConfInfo{
			ConfUri: in.ConfUri,
			Subject: in.Subject,
			Users:   in.Users,
		}
		err = initializer.DB.Create(&confInfo).Error
		if err != nil {
			return err
		}
	}

	fmt.Println("initialized mock data")
	return nil
}
