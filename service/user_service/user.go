package user_service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"my-go-blog/models"
)

var ErrInvalidPassword = errors.New("密码长度不合法")

func Register(username, password string) (models.User, error) {
	// 1. 检查密码的字节长度
	if len([]byte(password)) > 72 {
		return models.User{},
			ErrInvalidPassword
	}
	// 2. 检查用户名是否存在
	_, err := models.GetUserByUsername(username)
	if err == nil {
		return models.User{}, models.ErrUsernameExists
	}
	// 其他数据库错误不能当作用户名不存在
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, err
	}
	// 3. 将明文密码转换为 bcrypt 哈希
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return models.User{}, err
	}
	// 4. 调用 Model 层插入用户
	user, err := models.CreateUser(
		username,
		string(hashedPassword),
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
