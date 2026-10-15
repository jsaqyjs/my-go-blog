package models

import (
	"errors"
	"gorm.io/gorm"
	"my-go-blog/database"
)

var ErrUsernameExists = errors.New("用户名已存在")

type User struct {
	ID           uint   `gorm:"primary_key;auto_increment" json:"id"`
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
}

// TableName 指定用户表名称
func (User) TableName() string {
	return "users"
}

// GetUserByUsername 通过用户名
func GetUserByUsername(username string) (User, error) {
	var user User
	result := database.DB.Where("username=?", username).First(&user)
	return user, result.Error
}

// CreateUser 保存新用户
func CreateUser(username, passwordHash string) (User, error) {
	user := User{
		Username:     username,
		PasswordHash: passwordHash,
	}
	err := database.DB.Create(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return User{}, ErrUsernameExists
		}
		return User{}, err
	}
	return user, nil
}

// GetUserByID 根据主键查询用户
func GetUserByID(id uint) (User, error) {
	var user User
	err := database.DB.First(&user, id).Error
	return user, err
}
