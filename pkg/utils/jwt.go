package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const TokenExpiration = 7 * 24 * time.Hour

// HmacUser 签名需要传递的参数(根据自己定义)
type HmacUser struct {
	AccountId int64  `json:"accountId"` // 账号id
	Username  string `json:"username"`  // 用户名
	IsAdmin   int64  `json:"isAdmin"`   // 是否为超管1表是,2表示否
}

type MyClaims struct {
	AccountId int64  `json:"accountId"`
	Username  string `json:"username"`
	IsAdmin   int64  `json:"isAdmin"` // 是否为超管1表是,2表示否
	jwt.StandardClaims
}

// LoginStruct 登录的参数
type LoginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 证书签名密钥
var jwtKey = []byte("abc")

// GenerateToken 定义生成token的方法
func GenerateToken(u HmacUser) (string, error) {
	// 定义过期时间,7天后过期
	expirationTime := time.Now().Add(TokenExpiration)
	claims := &MyClaims{
		AccountId: u.AccountId,
		Username:  u.Username,
		IsAdmin:   u.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),     // 发布时间
			Subject:   "token",               // 主题
			Issuer:    "水痕",                  // 发布者
		},
	}
	// 注意单词别写错了
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthTokenRedisKey(accountId int64, token string) string {
	return fmt.Sprintf("auth:token:%d:%s", accountId, Md5(token))
}

// ParseToken 定义解析token的方法
func ParseToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
