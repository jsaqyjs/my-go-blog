package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

var secret = "my-go-blog"

func GenerateToken(userID uint) (string, error) {

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

func ParseToken(tokenString string) (uint, error) {
	// 2. 创建接收 JWT Payload 的结构体
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(secret),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt())
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("无效的 JWT")
	}
	// 4. 从 Subject 中读取用户 ID
	if claims.Subject == "" {
		return 0, errors.New("JWT 缺少用户 ID")
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	// 防止转换为 uint 时发生溢出
	if userID == 0 || userID > uint64(^uint(0)) {
		return 0, errors.New("无效的用户 ID")
	}

	return uint(userID), nil
}
