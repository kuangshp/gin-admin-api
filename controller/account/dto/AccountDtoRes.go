package dto

import (
	"fmt"
	"gin_admin_api/model"
	"time"
)

// 定义返回的数据模型
type AccountDtoRes struct {
	Id        uint      `json:"id"`
	UserName  string    `json:"username"`
	Mobile    string    `json:"mobile"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 将数据模型转换为返回值的
func ToAccountModelToRes(account model.Account) AccountDtoRes {
	return AccountDtoRes{
		Id:        account.ID,
		UserName:  account.UserName,
		Mobile:    account.Mobile,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}

// 列表的转换
func ToAccountModelListToRes(account []model.Account) []AccountDtoRes {
	result := make([]AccountDtoRes, 0)
	for _, item := range account {
		fmt.Println(item.UserName)
		result = append(result, AccountDtoRes{
			Id:        item.ID,
			UserName:  item.UserName,
			Mobile:    item.Mobile,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return result
}
