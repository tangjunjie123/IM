package models

import "github.com/jinzhu/gorm"

type UserBase struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	Identity      string
	ClintIp       string
	ClintPort     string
	LoginTime     uint64
	HeartbeatTime uint64
	LogoutTime    uint64
	IsLogout      bool
	DeviInfo      string
	Salt          string
}
