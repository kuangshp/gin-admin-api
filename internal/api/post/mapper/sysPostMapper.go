package mapper

import (
	"gin-admin-api/internal/api/post/dto"
	"gin-admin-api/internal/api/post/vo"
	"gin-admin-api/internal/dal/model"
)

type ISysPostMapper interface {
	DtoToEntity(req *dto.CreateSysPostDTO) *model.SysPostEntity
	EntityToVO(entity *model.SysPostEntity) vo.SysPostVO
}

type sysPostMapper struct{}

func NewSysPostMapper() ISysPostMapper {
	return &sysPostMapper{}
}

func (m *sysPostMapper) DtoToEntity(req *dto.CreateSysPostDTO) *model.SysPostEntity {
	return &model.SysPostEntity{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}
}

func (m *sysPostMapper) EntityToVO(entity *model.SysPostEntity) vo.SysPostVO {
	if entity == nil {
		return vo.SysPostVO{}
	}
	return vo.SysPostVO{
		ID:        entity.ID,
		Name:      entity.Name,
		Code:      entity.Code,
		Sort:      entity.Sort,
		Status:    entity.Status,
		Remark:    entity.Remark,
		CreatedAt: entity.CreatedAt.Unix(),
		UpdatedAt: entity.UpdatedAt.Unix(),
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
	}
}
