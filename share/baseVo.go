package share

import "gin-admin-api/model"

type BaseVo struct {
	ID        int64           `json:"id,string"`
	CreatedAt model.LocalTime `json:"createdAt"`
	UpdatedAt model.LocalTime `json:"updatedAt"`
}
