package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

//User DB schema
type User struct {
	gorm.Model

	Username      string `gorm:"type:varchar(100);unique_index;not null"`
	Name          string `gorm:"type:varchar(100);not null"`
	Lastname      string `gorm:"type:varchar(100);not null"`
	Password      string `gorm:"type:varchar(100);not null"`
	Email         string `gorm:"type:varchar(100);unique;not null"`
	Token         string `gorm:"type:varchar(255);unique;not null"`
	ApprovalToken string `gorm:"type:varchar(255);not null"`
	Approved      bool   `gorm:"type:boolean"`
}

//CreateUser create a new user
func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(&user).Error
}

//GetUser get a single user record
func GetUser(db *gorm.DB, username string) (User, error) {
	var user User

	if db.Where("username = ?", username).First(&user).RecordNotFound() {
		return user, errors.New("record not found")
	}
	return user, nil
}

//UpdateUser update an existing user
func UpdateUser(db *gorm.DB, user *User) error {
	return db.Save(&user).Error
}
