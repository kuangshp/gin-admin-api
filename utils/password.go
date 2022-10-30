package utils

import (
	"fmt"
	"gin-admin-api/global"
	"golang.org/x/crypto/bcrypt"
)

// MakePassword 对明文密码进行加密
func MakePassword(password string) (string, error) {
	salt := global.ServerConfig.Salt
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s", password, salt)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	} else {
		return string(hashPassword), nil
	}
}

// CheckPassword 校验密码是否正确
func CheckPassword(sqlPassword string, password string) (bool, error) {
	salt := global.ServerConfig.Salt
	err := bcrypt.CompareHashAndPassword([]byte(sqlPassword), []byte(fmt.Sprintf("%s%s", password, salt)))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
