package dto

import (
	"gin_admin_api/model"
)

type AccountDto struct {
	UserName string `json:"username" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=6,max=16"`
	Mobile   string `json:"mobile" binding:"min=11,max=11"`
}

func ToAccountDto(account model.Account) AccountDto {
	return AccountDto{
		UserName: account.UserName,
		Mobile:   account.Mobile,
	}
}

func ToAccountModel(account AccountDto) model.Account {
	return model.Account{
		UserName: account.UserName,
		Password: account.Password,
		Mobile:   account.Mobile,
	}
}
