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

// Input for mocking ConfInfos
type InputConfInfos struct {
	Data []ConfInfo
}

// Input for mocking ConfInfo
type InputConfInfo struct {
	ConfUri datatypes.URL
	Subject string
	Users   []User
}

// User in a conference
type User struct {
	UserID    uint `gorm:"primary_key"`
	EntityUrl datatypes.URL
	Role      string
	ConfRefer uint
}

// Conference Information
type ConfInfo struct {
	gorm.Model
	ConfUri datatypes.URL
	Subject string
	Users   []User `gorm:"foreignKey:ConfRefer"`
}

// GET /api/v1/confInfos
// Get all ConfInfos
func GetAllConfInfos(confInfos *[]ConfInfo) (err error) {
	if err = initializer.DB.
		Preload("Users").
		Find(&confInfos).
		Error; err != nil {
		return err
	}
	return nil
}

// GET /api/v1/confInfos/:id
// Get a ConfInfo by its id
func GetOneConfInfoById(confId string, confInfo *ConfInfo) (err error) {
	if err = initializer.DB.
		Where("id = ?", confId).
		Preload("Users").
		First(&confInfo).
		Error; err != nil {
		return err
	}
	return nil
}

// GET /api/v1/confInfos/:topic
// Get a ConfInfo by its topic which equals to its subject field
func GetOneConfInfoByTopic(topic string, confInfo *ConfInfo) (err error) {
	if err = initializer.DB.
		Where("subject = ?", topic).
		Preload("Users").
		First(&confInfo).
		Error; err != nil {
		return err
	}
	return nil
}

// POST /api/v1/confInfos
// Create a ConfInfo in
func AddNewConfInfo(confInfo *ConfInfo) (err error) {
	if err = initializer.DB.Create(&confInfo).Error; err != nil {
		return err
	}
	return nil
}

// PATCH /api/v1/topicMode/confInfos/:topic
// Patch a new User to a ConfInfo using a Topic
func PatchUserToTopic(topic string, newUser User, confInfo *ConfInfo) (err error) {
	if err = initializer.DB.
		Where("subject = ?", topic).
		First(&confInfo).
		Error; err != nil {
		return err
	}
	confInfo.Users = append(confInfo.Users, newUser)
	initializer.DB.Save(&confInfo)
	return nil
}

// DELETE /api/v1/topicMode/confInfos/:topic/users/:userUrl
// Delete a user from a ConfInfo using his userUrl
func DeleteUserFromTopic(topic string, user User, confInfo *ConfInfo) (err error) {
	if err = initializer.DB.Where("subject = ?", topic).Preload("Users").First(&confInfo).Error; err != nil {
		return err
	}
	if err = initializer.DB.Where("entity_url = ? AND role = ?", user.EntityUrl, user.Role).Delete(&confInfo.Users).Error; err != nil {
		return err
	}
	return nil
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
