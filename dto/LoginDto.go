package dto

type LoginDto struct {
	UserName string `validate:"required,checkName" json:"username"`
	Password string `validate:"required" json:"password"`
}
