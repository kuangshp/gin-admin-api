package dto

type AccountDto struct {
	Username string `json:"username" binding:"required,min=2,max=10"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type CreateAccountDto struct {
	AccountDto
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16"`
}

type ModifyAccountPassword struct {
	Password        string `json:"password" binding:"required,min=6,max=16"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16"`
}

type ModifyCurrentPassword struct {
	Password        string `json:"password" binding:"required,min=6,max=16"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=16"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16"`
}
