package mapper

import (
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/vo"
)

// ISysDeptMapper 部门表 mapper 接口
type ISysDeptMapper interface {
	// DtoToEntity 将 CreateSysDeptDTO 映射到 SysDeptEntity
	DtoToEntity(d *dto.CreateSysDeptDTO) *model.SysDeptEntity
	// EntityToVO 将 SysDeptEntity 映射到 SysDeptVO
	EntityToVO(e *model.SysDeptEntity) *vo.SysDeptVO
}

type sysDeptMapper struct{}

// NewSysDeptMapper 创建 SysDeptMapper 实例
func NewSysDeptMapper() ISysDeptMapper {
	return &sysDeptMapper{}
}

func (m *sysDeptMapper) DtoToEntity(d *dto.CreateSysDeptDTO) *model.SysDeptEntity {
	return &model.SysDeptEntity{
		Name:     d.Name,
		ParentID: d.ParentID,
		FullID:   d.FullID,
		FullName: d.FullName,
		Sort:     d.Sort,
		Status:   d.Status,
		LeaderID: d.LeaderID,
		Phone:    d.Phone,
		Email:    d.Email,
	}
}

func (m *sysDeptMapper) EntityToVO(e *model.SysDeptEntity) *vo.SysDeptVO {
	if e == nil {
		return nil
	}
	return &vo.SysDeptVO{
		ID:        e.ID,
		Name:      e.Name,
		ParentID:  e.ParentID,
		FullID:    e.FullID,
		FullName:  e.FullName,
		Sort:      e.Sort,
		Status:    e.Status,
		LeaderID:  e.LeaderID,
		Phone:     e.Phone,
		Email:     e.Email,
		CreatedAt: e.CreatedAt.Unix(),
		UpdatedAt: e.UpdatedAt.Unix(),
		CreatedBy: e.CreatedBy,
		UpdatedBy: e.UpdatedBy,
	}
}
