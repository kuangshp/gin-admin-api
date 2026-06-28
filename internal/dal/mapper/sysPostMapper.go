package mapper

import (
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/vo"
)

// ISysPostMapper 岗位表 mapper 接口
type ISysPostMapper interface {
	// DtoToEntity 将 CreateSysPostDTO 映射到 SysPostEntity
	DtoToEntity(d *dto.CreateSysPostDTO) *model.SysPostEntity
	// EntityToVO 将 SysPostEntity 映射到 SysPostVO
	EntityToVO(e *model.SysPostEntity) *vo.SysPostVO
}

type sysPostMapper struct{}

// NewSysPostMapper 创建 SysPostMapper 实例
func NewSysPostMapper() ISysPostMapper {
	return &sysPostMapper{}
}

func (m *sysPostMapper) DtoToEntity(d *dto.CreateSysPostDTO) *model.SysPostEntity {
	return &model.SysPostEntity{
		Name:   d.Name,
		Code:   d.Code,
		Sort:   d.Sort,
		Status: d.Status,
		Remark: d.Remark,
	}
}

func (m *sysPostMapper) EntityToVO(e *model.SysPostEntity) *vo.SysPostVO {
	if e == nil {
		return nil
	}
	return &vo.SysPostVO{
		ID:        e.ID,
		Name:      e.Name,
		Code:      e.Code,
		Sort:      e.Sort,
		Status:    e.Status,
		Remark:    e.Remark,
		CreatedAt: e.CreatedAt.Unix(),
		UpdatedAt: e.UpdatedAt.Unix(),
		CreatedBy: e.CreatedBy,
		UpdatedBy: e.UpdatedBy,
	}
}
