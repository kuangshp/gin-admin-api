package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// MakePassword 对明文密码进行加密
func MakePassword(password, salt string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s", password, salt)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	} else {
		return string(hashPassword), nil
	}
}

// CheckPassword 校验密码是否正确
func CheckPassword(sqlPassword, password, salt string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(sqlPassword), []byte(fmt.Sprintf("%s%s", password, salt)))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
