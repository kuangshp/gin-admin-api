package dept

import (
	"errors"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/api/dept/dto"
	"gin-admin-api/internal/api/dept/mapper"
	"gin-admin-api/internal/api/dept/vo"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuangshp/go-utils/k"
	"github.com/kuangshp/gorm-plus"
	"github.com/spf13/cast"
	"gorm.io/gen"
)

type IDept interface {
	CreateDeptApi(ctx *gin.Context)
	DeleteDeptByIdApi(ctx *gin.Context)
	ModifyDeptByIdApi(ctx *gin.Context)
	GetDeptPageApi(ctx *gin.Context)
	GetDeptListApi(ctx *gin.Context)
	GetDeptDetailByIdApi(ctx *gin.Context)
}

type Dept struct {
	*base.BaseApi
	DeptRepository    repository.SysDeptRepository
	AccountRepository repository.SysAccountRepository
	DeptMapper        mapper.ISysDeptMapper
}

// CreateDeptApi 创建部门
// @Summary 创建部门
// @Description 创建部门；部门层级 fullId/fullName 由后端根据 parentId 自动维护
// @Tags 部门中心
// @Accept json
// @Produce json
// @Param data body dto.CreateSysDeptDTO true "创建部门参数"
// @Success 200 {string} string "创建成功"
// @Router /api/v1/admin/dept [post]
func (d Dept) CreateDeptApi(ctx *gin.Context) {
	var req dto.CreateSysDeptDTO
	if !d.BindAndValidateJSON(ctx, &req) {
		return
	}
	if req.Status == 0 {
		req.Status = enum.StatusNormalEnum
	}
	if req.ParentID == 0 && ctx.GetInt64("isAdmin") != 1 {
		d.Fail(ctx, errors.New("无权限创建根部门"), "无权限创建根部门")
		return
	}
	parent, fullID, fullName, ok := d.buildDeptPath(ctx, req.ParentID, req.Name, 0)
	if !ok {
		return
	}
	if !d.checkNameUnique(ctx, req.ParentID, req.Name, 0) {
		return
	}
	if !d.checkLeaderValid(ctx, req.LeaderID) {
		return
	}
	if parent == nil && req.ParentID != 0 {
		d.Fail(ctx, errors.New("父部门不存在"), "父部门不存在或无权限访问")
		return
	}
	entity := d.DeptMapper.DtoToEntity(&req, fullID, fullName)
	if err := d.DeptRepository.Create(ctx, entity); err != nil {
		d.Fail(ctx, err, "创建部门失败")
		return
	}
	d.Success(ctx, "创建成功")
	return
}

// DeleteDeptByIdApi 删除部门
// @Summary 删除部门
// @Description 根据部门 ID 删除部门；存在子部门或账号时不能删除
// @Tags 部门中心
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {string} string "删除成功"
// @Router /api/v1/admin/dept/{id} [delete]
func (d Dept) DeleteDeptByIdApi(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	if id <= 0 {
		d.Fail(ctx, errors.New("参数id错误"), "参数id错误")
		return
	}
	if _, err := d.DeptRepository.FindById(ctx, id); err != nil {
		d.Fail(ctx, err, "部门不存在或无权限访问")
		return
	}
	existsChild, err := d.DeptRepository.Exists(ctx, gormplus.QueryOpt().Where(dao.SysDeptEntity.ParentID.Eq(id)).Build())
	if err != nil {
		d.Fail(ctx, err, "查询子部门失败")
		return
	}
	if existsChild {
		d.Fail(ctx, errors.New("存在子部门"), "存在子部门，不能删除")
		return
	}
	existsAccount, err := d.AccountRepository.Exists(ctx, gormplus.QueryOpt().Where(dao.SysAccountEntity.DeptID.Eq(id)).Build())
	if err != nil {
		d.Fail(ctx, err, "查询账号失败")
		return
	}
	if existsAccount {
		d.Fail(ctx, errors.New("部门下存在账号"), "部门下存在账号，不能删除")
		return
	}
	if err := d.DeptRepository.DeleteById(ctx, id); err != nil {
		d.Fail(ctx, err, "删除部门失败")
		return
	}
	d.Success(ctx, "删除成功")
	return
}

// ModifyDeptByIdApi 修改部门
// @Summary 修改部门
// @Description 根据部门 ID 修改部门基础信息和父级关系
// @Tags 部门中心
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Param data body dto.ModifySysDeptDTO true "修改部门参数"
// @Success 200 {string} string "修改成功"
// @Router /api/v1/admin/dept/{id} [put]
func (d Dept) ModifyDeptByIdApi(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	if id <= 0 {
		d.Fail(ctx, errors.New("参数id错误"), "参数id错误")
		return
	}
	var req dto.CreateSysDeptDTO
	if !d.BindAndValidateJSON(ctx, &req) {
		return
	}
	if req.ParentID == id {
		d.Fail(ctx, errors.New("上级部门不能是自己"), "上级部门不能是自己")
		return
	}
	if _, err := d.DeptRepository.FindById(ctx, id); err != nil {
		d.Fail(ctx, err, "部门不存在或无权限访问")
		return
	}
	if req.Status == 0 {
		req.Status = enum.StatusNormalEnum
	}
	_, fullID, fullName, ok := d.buildDeptPath(ctx, req.ParentID, req.Name, id)
	if !ok {
		return
	}
	if !d.checkNameUnique(ctx, req.ParentID, req.Name, id) {
		return
	}
	if !d.checkLeaderValid(ctx, req.LeaderID) {
		return
	}
	if err := d.DeptRepository.UpdateById(ctx, id, gormplus.Update().WithColumns(
		dao.SysDeptEntity.Name.Value(req.Name),
		dao.SysDeptEntity.ParentID.Value(req.ParentID),
		dao.SysDeptEntity.FullID.Value(fullID),
		dao.SysDeptEntity.FullName.Value(fullName),
		dao.SysDeptEntity.Sort.Value(req.Sort),
		dao.SysDeptEntity.Status.Value(req.Status),
		dao.SysDeptEntity.LeaderID.Value(req.LeaderID),
		dao.SysDeptEntity.Phone.Value(req.Phone),
		dao.SysDeptEntity.Email.Value(req.Email),
	).Build()); err != nil {
		d.Fail(ctx, err, "修改部门失败")
		return
	}
	d.Success(ctx, "修改成功")
	return
}

// GetDeptPageApi 分页获取部门
// @Summary 分页获取部门
// @Description 根据部门名称、状态、上级部门分页查询部门列表，自动应用数据权限
// @Tags 部门中心
// @Accept json
// @Produce json
// @Param data body dto.SysDeptPageDTO true "部门分页查询参数"
// @Success 200 {array} vo.SysDeptVO "部门分页列表"
// @Router /api/v1/admin/dept/page [post]
func (d Dept) GetDeptPageApi(ctx *gin.Context) {
	var req dto.SysDeptPageDTO
	if !d.BindAndValidateJSON(ctx, &req) {
		return
	}
	list, total, err := d.DeptRepository.FindPageByWrapper(ctx, req.PageNumber, req.PageSize, func(g gormplus.IGenWrapper[dao.ISysDeptEntityDo]) {
		g.WhereIf(req.Name != "", dao.SysDeptEntity.Name.Like("%"+req.Name+"%")).
			WhereIf(req.Status > 0, dao.SysDeptEntity.Status.Eq(req.Status)).
			WhereIf(req.ParentID >= 0, dao.SysDeptEntity.ParentID.Eq(req.ParentID))
	}, gormplus.QueryOpt().Order(dao.SysDeptEntity.Sort, dao.SysDeptEntity.ID).Build())
	if err != nil {
		d.Fail(ctx, err, "获取部门分页失败")
		return
	}
	d.BuildPageData(ctx, k.Map(list, func(item *model.SysDeptEntity, index int) vo.SysDeptVO {
		return d.DeptMapper.EntityToVO(item)
	}), total)
}

// GetDeptListApi 获取部门列表
// @Summary 获取部门列表
// @Description 根据部门名称和状态获取部门列表，自动应用数据权限
// @Tags 部门中心
// @Accept json
// @Produce json
// @Param name query string false "部门名称"
// @Param status query int false "状态：1正常，2禁用，0全部"
// @Success 200 {array} vo.SysDeptVO "部门列表"
// @Router /api/v1/admin/dept/list [get]
func (d Dept) GetDeptListApi(ctx *gin.Context) {
	name := ctx.DefaultQuery("name", "")
	status := cast.ToInt64(ctx.DefaultQuery("status", "0"))
	list, err := d.DeptRepository.FindListByWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysDeptEntityDo]) {
		g.WhereIf(name != "", dao.SysDeptEntity.Name.Like("%"+name+"%")).
			WhereIf(status > 0, dao.SysDeptEntity.Status.Eq(status))
	}, gormplus.QueryOpt().Order(dao.SysDeptEntity.Sort, dao.SysDeptEntity.ID).Build())
	if err != nil {
		d.Fail(ctx, err, "获取部门列表失败")
		return
	}
	d.Success(ctx, k.Map(list, func(item *model.SysDeptEntity, index int) vo.SysDeptVO {
		return d.DeptMapper.EntityToVO(item)
	}))
}

// GetDeptDetailByIdApi 获取部门详情
// @Summary 获取部门详情
// @Description 根据部门 ID 获取部门详情，自动应用数据权限
// @Tags 部门中心
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} vo.SysDeptVO "部门详情"
// @Router /api/v1/admin/dept/detail/{id} [get]
func (d Dept) GetDeptDetailByIdApi(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	if id <= 0 {
		d.Fail(ctx, errors.New("参数id错误"), "参数id错误")
		return
	}
	entity, err := d.DeptRepository.FindById(ctx, id)
	if err != nil {
		d.Fail(ctx, err, "部门不存在或无权限访问")
		return
	}
	d.Success(ctx, d.DeptMapper.EntityToVO(entity))
}

// 组装部门数据
func (d Dept) buildDeptPath(ctx *gin.Context, parentID int64, name string, selfID int64) (*model.SysDeptEntity, string, string, bool) {
	if parentID == 0 {
		return nil, "", name, true
	}
	parent, err := d.DeptRepository.FindById(ctx, parentID)
	if err != nil {
		d.Fail(ctx, err, "父部门不存在或无权限访问")
		return nil, "", "", false
	}
	if selfID > 0 && containsDeptID(parent.FullID, selfID) {
		d.Fail(ctx, errors.New("不能移动到自己的下级部门"), "不能移动到自己的下级部门")
		return nil, "", "", false
	}
	fullID := strconv.FormatInt(parent.ID, 10)
	if parent.FullID != "" {
		fullID = parent.FullID + "," + strconv.FormatInt(parent.ID, 10)
	}
	fullName := name
	if parent.FullName != "" {
		fullName = parent.FullName + "/" + name
	}
	return parent, fullID, fullName, true
}

// 判断唯一性
func (d Dept) checkNameUnique(ctx *gin.Context, parentID int64, name string, excludeID int64) bool {
	conds := []gen.Condition{
		dao.SysDeptEntity.ParentID.Eq(parentID),
		dao.SysDeptEntity.Name.Eq(name),
	}
	if excludeID > 0 {
		conds = append(conds, dao.SysDeptEntity.ID.Neq(excludeID))
	}
	exists, err := d.DeptRepository.Exists(ctx, gormplus.QueryOpt().Where(conds...).Build())
	if err != nil {
		d.Fail(ctx, err, "部门名称重复校验失败")
		return false
	}
	if exists {
		d.Fail(ctx, errors.New("同级部门名称已存在"), "同级部门名称已存在")
		return false
	}
	return true
}

// 判断部门负责人
func (d Dept) checkLeaderValid(ctx *gin.Context, leaderID int64) bool {
	if leaderID == 0 {
		return true
	}
	exists, err := d.AccountRepository.Exists(ctx, gormplus.QueryOpt().Where(
		dao.SysAccountEntity.ID.Eq(leaderID),
		dao.SysAccountEntity.Status.Eq(enum.StatusNormalEnum),
	).Build())
	if err != nil {
		d.Fail(ctx, err, "负责人账号校验失败")
		return false
	}
	if !exists {
		d.Fail(ctx, errors.New("负责人账号不存在或无权限访问"), "负责人账号不存在或无权限访问")
		return false
	}
	return true
}

func containsDeptID(fullID string, deptID int64) bool {
	if fullID == "" || deptID <= 0 {
		return false
	}
	target := strconv.FormatInt(deptID, 10)
	for _, id := range strings.Split(fullID, ",") {
		if id == target {
			return true
		}
	}
	return false
}

func NewDept(baseApi *base.BaseApi) IDept {
	return Dept{
		BaseApi:           baseApi,
		DeptRepository:    repository.NewSysDeptRepository(),
		AccountRepository: repository.NewSysAccountRepository(),
		DeptMapper:        mapper.NewSysDeptMapper(),
	}
}
