package model

import (
	"gorm.io/datatypes"
)

type User struct {
	Entity datatypes.URL
	Role   string
}
