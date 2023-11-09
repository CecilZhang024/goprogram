package modules

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Username string `gorm:"unique"`
	Password string
	Type     string
}
