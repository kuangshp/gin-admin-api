package dto

// CreateCasbinRuleDTO 创建CasbinRule请求
type CreateCasbinRuleDTO struct {
	Ptype string `json:"ptype" validate:"omitempty,max=100"` // 
	V0 string `json:"v0" validate:"omitempty,max=100"` // 
	V1 string `json:"v1" validate:"omitempty,max=100"` // 
	V2 string `json:"v2" validate:"omitempty,max=100"` // 
	V3 string `json:"v3" validate:"omitempty,max=100"` // 
	V4 string `json:"v4" validate:"omitempty,max=100"` // 
	V5 string `json:"v5" validate:"omitempty,max=100"` // 
}

// ModifyCasbinRuleDTO 修改CasbinRule请求
type ModifyCasbinRuleDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateCasbinRuleDTO
}
