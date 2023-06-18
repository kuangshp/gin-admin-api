package share

import "gin-admin-api/model"

type BaseVo struct {
	ID        int64           `json:"id,string"`
	CreatedAt model.LocalTime `json:"createdAt"`
	UpdatedAt model.LocalTime `json:"updatedAt"`
}

type PageDataVo struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}
