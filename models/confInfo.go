package model

import (
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

// POST /confInfos
// Create a ConfInfo in
func AddNewConfInfo(confInfo *ConfInfo) (err error) {
	if err = initializer.DB.Create(&confInfo).Error; err != nil {
		return err
	}
	return nil
}
