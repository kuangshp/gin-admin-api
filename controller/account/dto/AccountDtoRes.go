package dto

import (
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
//func ToAccountModelListToRes(account []model.Account) {
//	result := make([]AccountDtoRes, 10)
//	for key, val := range account {
//		result = append(result, {key,val})
//	}
//}
