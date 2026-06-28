package mapper

import (
	"gin-admin-api/internal/api/dept/dto"
	"gin-admin-api/internal/api/dept/vo"
	"gin-admin-api/internal/dal/model"
)

type ISysDeptMapper interface {
	DtoToEntity(req *dto.CreateSysDeptDTO, fullID, fullName string) *model.SysDeptEntity
	EntityToVO(entity *model.SysDeptEntity) vo.SysDeptVO
}

type sysDeptMapper struct{}

func NewSysDeptMapper() ISysDeptMapper {
	return &sysDeptMapper{}
}

func (m *sysDeptMapper) DtoToEntity(req *dto.CreateSysDeptDTO, fullID, fullName string) *model.SysDeptEntity {
	return &model.SysDeptEntity{
		Name:     req.Name,
		ParentID: req.ParentID,
		FullID:   fullID,
		FullName: fullName,
		Sort:     req.Sort,
		Status:   req.Status,
		LeaderID: req.LeaderID,
		Phone:    req.Phone,
		Email:    req.Email,
	}
}

func (m *sysDeptMapper) EntityToVO(entity *model.SysDeptEntity) vo.SysDeptVO {
	if entity == nil {
		return vo.SysDeptVO{}
	}
	return vo.SysDeptVO{
		ID:        entity.ID,
		Name:      entity.Name,
		ParentID:  entity.ParentID,
		FullID:    entity.FullID,
		FullName:  entity.FullName,
		Sort:      entity.Sort,
		Status:    entity.Status,
		LeaderID:  entity.LeaderID,
		Phone:     entity.Phone,
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt.Unix(),
		UpdatedAt: entity.UpdatedAt.Unix(),
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
	}
}
