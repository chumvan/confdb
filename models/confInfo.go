package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ConfInfo struct {
	gorm.Model
	ConfId  string
	ConfUri datatypes.URL
	ConfDesc
	ConfState
}
