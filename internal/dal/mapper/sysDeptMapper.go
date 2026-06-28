package mapper

import (
	"time"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/vo"
)

// ISysDeptMapper 部门表 mapper 接口
type ISysDeptMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysDeptDTO) *model.SysDeptEntity
	// EntityToVO 将数据库实体映射到响应结构体
	EntityToVO(e *model.SysDeptEntity) *vo.SysDeptVO
}

// sysDeptMapper mapper 实现
type sysDeptMapper struct{}

// NewSysDeptMapper 创建 SysDeptMapper 实例
func NewSysDeptMapper() ISysDeptMapper {
	return &sysDeptMapper{}
}

// DtoToEntity 将 CreateSysDeptDTO 映射到 SysDeptEntity
func (m *sysDeptMapper) DtoToEntity(d *dto.CreateSysDeptDTO) *model.SysDeptEntity {
	e := &model.SysDeptEntity{
		Name: d.Name, // 部门名称
		ParentID: d.ParentID, // 上级部门id
		FullID: d.FullID, // 全层级ID，例：1,5,12
		FullName: d.FullName, // 全层级名称
		Sort: d.Sort, // 排序
		Status: d.Status, // 状态1正常 2禁用
		LeaderID: d.LeaderID, // 负责人账号id，关联sys_account.id
		Phone: d.Phone, // 联系电话
		Email: d.Email, // 邮箱
	}
	return e
}

// EntityToVO 将 SysDeptEntity 映射到 SysDeptVO
func (m *sysDeptMapper) EntityToVO(e *model.SysDeptEntity) *vo.SysDeptVO {
	if e == nil {
		return nil
	}
	return &vo.SysDeptVO{
		ID: e.ID, // 主键id
		Name: e.Name, // 部门名称
		ParentID: e.ParentID, // 上级部门id
		FullID: e.FullID, // 全层级ID，例：1,5,12
		FullName: e.FullName, // 全层级名称
		Sort: e.Sort, // 排序
		Status: e.Status, // 状态1正常 2禁用
		LeaderID: e.LeaderID, // 负责人账号id，关联sys_account.id
		Phone: e.Phone, // 联系电话
		Email: e.Email, // 邮箱
        CreatedAt: e.CreatedAt.Unix(), // 创建时间
        UpdatedAt: e.UpdatedAt.Unix(), // 更新时间
		CreatedBy: e.CreatedBy, // 创建人
		UpdatedBy: e.UpdatedBy, // 更新人

	}
}
