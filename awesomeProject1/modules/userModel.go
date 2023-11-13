package modules

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Username string `gorm:"unique"`
	Password string
	Type     string
}
type UserAuth struct {
	Token string
}
type Userinfo struct {
	Roles        [1]string
	Introduction string
	Avatar       string
	Name         string
}
