package vo

import (
	"fmt"
	"gin_admin_api/model"
	"time"
)

// AccountVo 定义返回的数据模型
type AccountVo struct {
	Id       int32  `json:"id"`
	Address  string `json:"address"`
	Avatar   string `json:"avatar"`
	Desc     string `json:"desc"`
	Gender   string `json:"gender"`
	BirthDay *time.Time `json:"birth_day"`
	//UserName string `json:"username"`
	//Mobile   string `json:"mobile"`
	//CreatedAt time.Time `json:"created_at"`
	//UpdatedAt time.Time `json:"updated_at"`
	// 时间转换回去
	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at"`
}

// ToAccountModelToRes 将数据模型转换为返回值的
func ToAccountModelToRes(account model.AccountEntity) AccountVo {
	return AccountVo{
		Id: account.ID,
		//UserName:  account.UserName,
		//Mobile:    account.Mobile,
		Address:   account.Address,
		Avatar:    account.Avatar,
		Desc:      account.Desc,
		Gender:    account.Gender,
		//BirthDay:  uint64(account.BirthDay.Unix()),
		BirthDay:  account.BirthDay,
		CreatedAt: uint64(account.CreatedAt.Unix()),
		UpdatedAt: uint64(account.UpdatedAt.Unix()),
	}
}

// ToAccountModelListToRes 列表的转换
func ToAccountModelListToRes(account []model.AccountEntity) []AccountVo {
	result := make([]AccountVo, 0)
	for _, item := range account {
		fmt.Println(item.UserName)
		result = append(result, ToAccountModelToRes(item))
	}
	return result
}
