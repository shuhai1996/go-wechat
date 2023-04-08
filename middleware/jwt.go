package middleware

import "C"
import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"go-wechat/config"
)

// CustomClaims 自定义声明
type CustomClaims struct {
	UUID     uuid.UUID `json:"uuid"`
	ID       int `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nick_name"`
	jwt.StandardClaims
}

type User struct {
	ID       int `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
// JwtAuth 后台鉴权中间件
func JwtAuth(jwtS string) (user *User, err error){
	var secret = config.Get("jwt.secret")
	if jwtS == "" {
		err = fmt.Errorf("empty token")
		return
	}

	// 解析Token
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(jwtS, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// 验证Token是否有效
	if !token.Valid || err!=nil{
		return
	}

	if err != nil {
		return
	}

	return &User{
		Username: claims.Username,
		Nickname: claims.NickName,
		ID: claims.ID,
	}, nil
}
