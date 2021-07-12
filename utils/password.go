package utils

import "golang.org/x/crypto/bcrypt"

// GeneratePassword 对明文密码进行加密
func GeneratePassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	} else {
		return string(hashPassword), nil
	}
}

// CheckPassword 校验密码是否正确
func CheckPassword(sqlPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(sqlPassword), []byte(password))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
