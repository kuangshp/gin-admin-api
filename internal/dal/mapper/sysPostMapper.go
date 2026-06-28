package mapper

import (
	"time"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/vo"
)

// ISysPostMapper 岗位表 mapper 接口
type ISysPostMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysPostDTO) *model.SysPostEntity
	// EntityToVO 将数据库实体映射到响应结构体
	EntityToVO(e *model.SysPostEntity) *vo.SysPostVO
}

// sysPostMapper mapper 实现
type sysPostMapper struct{}

// NewSysPostMapper 创建 SysPostMapper 实例
func NewSysPostMapper() ISysPostMapper {
	return &sysPostMapper{}
}

// DtoToEntity 将 CreateSysPostDTO 映射到 SysPostEntity
func (m *sysPostMapper) DtoToEntity(d *dto.CreateSysPostDTO) *model.SysPostEntity {
	e := &model.SysPostEntity{
		Name: d.Name, // 岗位名称
		Code: d.Code, // 岗位编码
		Sort: d.Sort, // 排序
		Status: d.Status, // 状态1正常 2禁用
		Remark: d.Remark, // 备注
	}
	return e
}

// EntityToVO 将 SysPostEntity 映射到 SysPostVO
func (m *sysPostMapper) EntityToVO(e *model.SysPostEntity) *vo.SysPostVO {
	if e == nil {
		return nil
	}
	return &vo.SysPostVO{
		ID: e.ID, // 主键id
		Name: e.Name, // 岗位名称
		Code: e.Code, // 岗位编码
		Sort: e.Sort, // 排序
		Status: e.Status, // 状态1正常 2禁用
		Remark: e.Remark, // 备注
        CreatedAt: e.CreatedAt.Unix(), // 创建时间
        UpdatedAt: e.UpdatedAt.Unix(), // 更新时间
		CreatedBy: e.CreatedBy, // 创建人
		UpdatedBy: e.UpdatedBy, // 更新人

	}
}
