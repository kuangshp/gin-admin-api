package dto

type RegisterDto struct {
	UserName string `binding:"required,checkName" json:"username"`
	Password string `binding:"required" json:"password"`
}

