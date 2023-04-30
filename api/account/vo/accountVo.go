package vo

import (
	"gin-admin-api/model"
)

type AccountVo struct {
	Id        int64           `json:"id"`
	CreatedAt model.LocalTime `json:"createdAt"`
	UpdatedAt model.LocalTime `json:"updatedAt"`
	Username  string          `json:"username"` // 用户名
	Status    int64           `json:"status"`   // 状态1是正常,0是禁用
}
