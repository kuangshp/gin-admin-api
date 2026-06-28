package mapper

import (
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/vo"
)

// ISysAccountPostMapper 账号岗位关联表 mapper 接口
type ISysAccountPostMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysAccountPostDTO) *model.SysAccountPostEntity
	// EntityToVO 将数据库实体映射到响应结构体
	EntityToVO(e *model.SysAccountPostEntity) *vo.SysAccountPostVO
}

// sysAccountPostMapper mapper 实现
type sysAccountPostMapper struct{}

// NewSysAccountPostMapper 创建 SysAccountPostMapper 实例
func NewSysAccountPostMapper() ISysAccountPostMapper {
	return &sysAccountPostMapper{}
}

// DtoToEntity 将 CreateSysAccountPostDTO 映射到 SysAccountPostEntity
func (m *sysAccountPostMapper) DtoToEntity(d *dto.CreateSysAccountPostDTO) *model.SysAccountPostEntity {
	e := &model.SysAccountPostEntity{
		AccountID: d.AccountID, // 账号id
		PostID:    d.PostID,    // 岗位id
		IsPrimary: d.IsPrimary, // 是否主岗,1表示不是,2表示是
	}
	return e
}

// EntityToVO 将 SysAccountPostEntity 映射到 SysAccountPostVO
func (m *sysAccountPostMapper) EntityToVO(e *model.SysAccountPostEntity) *vo.SysAccountPostVO {
	if e == nil {
		return nil
	}
	return &vo.SysAccountPostVO{
		ID:        e.ID,               // 主键id
		AccountID: e.AccountID,        // 账号id
		PostID:    e.PostID,           // 岗位id
		IsPrimary: e.IsPrimary,        // 是否主岗,1表示不是,2表示是
		CreatedAt: e.CreatedAt.Unix(), // 创建时间
		UpdatedAt: e.UpdatedAt.Unix(), // 更新时间
		CreatedBy: e.CreatedBy,        // 创建人
		UpdatedBy: e.UpdatedBy,        // 更新人

	}
}
