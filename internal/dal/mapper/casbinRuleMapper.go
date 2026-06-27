package mapper

import (
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/vo"
)

// ICasbinRuleMapper CasbinRule mapper 接口
type ICasbinRuleMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateCasbinRuleDTO) *model.CasbinRuleEntity
	// EntityToVO 将数据库实体映射到响应结构体
	EntityToVO(e *model.CasbinRuleEntity) *vo.CasbinRuleVO
}

// casbinRuleMapper mapper 实现
type casbinRuleMapper struct{}

// NewCasbinRuleMapper 创建 CasbinRuleMapper 实例
func NewCasbinRuleMapper() ICasbinRuleMapper {
	return &casbinRuleMapper{}
}

// DtoToEntity 将 CreateCasbinRuleDTO 映射到 CasbinRuleEntity
func (m *casbinRuleMapper) DtoToEntity(d *dto.CreateCasbinRuleDTO) *model.CasbinRuleEntity {
	e := &model.CasbinRuleEntity{
		Ptype: d.Ptype, 
		V0: d.V0, 
		V1: d.V1, 
		V2: d.V2, 
		V3: d.V3, 
		V4: d.V4, 
		V5: d.V5, 
	}
	return e
}

// EntityToVO 将 CasbinRuleEntity 映射到 CasbinRuleVO
func (m *casbinRuleMapper) EntityToVO(e *model.CasbinRuleEntity) *vo.CasbinRuleVO {
	if e == nil {
		return nil
	}
	return &vo.CasbinRuleVO{
		ID: e.ID, 
		Ptype: e.Ptype, 
		V0: e.V0, 
		V1: e.V1, 
		V2: e.V2, 
		V3: e.V3, 
		V4: e.V4, 
		V5: e.V5, 

	}
}
