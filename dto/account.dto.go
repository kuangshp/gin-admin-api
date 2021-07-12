package dto

type AccountDto struct {
	UserName string `json:"username" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=6,max=16"`
	Mobile   string `json:"mobile" binding:"min=11,max=11"`
}
