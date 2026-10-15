package user_service

import (
	"errors"
	"my-go-blog/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"my-go-blog/models"
)

var ErrInvalidPassword = errors.New("密码长度不合法")
var ErrInvalidCredentials = errors.New("用户名或密码错误")

func Login(username string, password string) (string, error) {
	// 1. 根据用户名查询用户
	user, err := models.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	// 检查密码
	// 拿旧密码的加盐+cost 进行完全加密
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	// 3. 密码正确，生成 JWT
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetCurrentUser 查询当前登录用户
func GetCurrentUser(userID uint) (models.User, error) {
	return models.GetUserByID(userID)
}
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
