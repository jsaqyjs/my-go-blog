package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

func GenerateToken(userID uint) (string, error) {
	secret := "my-go-blog"
	now := time.Now()
	// claims 注册 jwt 带的一些参数
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
		Issuer:    "my-go-blog",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
