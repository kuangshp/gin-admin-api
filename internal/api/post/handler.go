package post

import (
	"errors"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/api/post/dto"
	"gin-admin-api/internal/api/post/mapper"
	"gin-admin-api/internal/api/post/vo"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"

	"github.com/gin-gonic/gin"
	"github.com/kuangshp/go-utils/k"
	"github.com/kuangshp/gorm-plus"
	"github.com/spf13/cast"
	"gorm.io/gen"
)

type IPost interface {
	CreatePostApi(ctx *gin.Context)
	DeletePostByIdApi(ctx *gin.Context)
	ModifyPostByIdApi(ctx *gin.Context)
	GetPostPageApi(ctx *gin.Context)
	GetPostListApi(ctx *gin.Context)
	GetPostDetailByIdApi(ctx *gin.Context)
}

type Post struct {
	*base.BaseApi
	PostRepository        repository.SysPostRepository
	AccountPostRepository repository.SysAccountPostRepository
	PostMapper            mapper.ISysPostMapper
}

// CreatePostApi 创建岗位
// @Summary 创建岗位
// @Description 创建岗位，岗位编码全局唯一
// @Tags 岗位中心
// @Accept json
// @Produce json
// @Param data body dto.CreateSysPostDTO true "创建岗位参数"
// @Success 200 {string} string "创建成功"
// @Router /api/v1/admin/post [post]
func (p Post) CreatePostApi(ctx *gin.Context) {
	var req dto.CreateSysPostDTO
	if !p.BindAndValidateJSON(ctx, &req) {
		return
	}
	if req.Status == 0 {
		req.Status = enum.StatusNormalEnum
	}
	if !p.checkCodeUnique(ctx, req.Code, 0) {
		return
	}
	if err := p.PostRepository.Create(ctx, p.PostMapper.DtoToEntity(&req)); err != nil {
		p.Fail(ctx, err, "创建岗位失败")
		return
	}
	p.Success(ctx, "创建成功")
}

// DeletePostByIdApi 删除岗位
// @Summary 删除岗位
// @Description 根据岗位 ID 删除岗位；已绑定账号时不能删除
// @Tags 岗位中心
// @Accept json
// @Produce json
// @Param id path int true "岗位ID"
// @Success 200 {string} string "删除成功"
// @Router /api/v1/admin/post/{id} [delete]
func (p Post) DeletePostByIdApi(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	if id <= 0 {
		p.Fail(ctx, errors.New("参数id错误"), "参数id错误")
		return
	}
	if _, err := p.PostRepository.FindById(ctx, id); err != nil {
		p.Fail(ctx, err, "岗位不存在或无权限访问")
		return
	}
	exists, err := p.AccountPostRepository.Exists(ctx, gormplus.QueryOpt().Where(dao.SysAccountPostEntity.PostID.Eq(id)).Build())
	if err != nil {
		p.Fail(ctx, err, "查询账号岗位失败")
		return
	}
	if exists {
		p.Fail(ctx, errors.New("岗位已绑定账号"), "岗位已绑定账号，不能删除")
		return
	}
	if err := p.PostRepository.DeleteById(ctx, id); err != nil {
		p.Fail(ctx, err, "删除岗位失败")
		return
	}
	p.Success(ctx, "删除成功")
}

// ModifyPostByIdApi 修改岗位
// @Summary 修改岗位
// @Description 根据岗位 ID 修改岗位基础信息
// @Tags 岗位中心
// @Accept json
// @Produce json
// @Param id path int true "岗位ID"
// @Param data body dto.ModifySysPostDTO true "修改岗位参数"
// @Success 200 {string} string "修改成功"
// @Router /api/v1/admin/post/{id} [put]
func (p Post) ModifyPostByIdApi(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	if id <= 0 {
		p.Fail(ctx, errors.New("参数id错误"), "参数id错误")
		return
	}
	var req dto.ModifySysPostDTO
	if !p.BindAndValidateJSON(ctx, &req) {
		return
	}
	if req.Status == 0 {
		req.Status = enum.StatusNormalEnum
	}
	if _, err := p.PostRepository.FindById(ctx, id); err != nil {
		p.Fail(ctx, err, "岗位不存在或无权限访问")
		return
	}
	if !p.checkCodeUnique(ctx, req.Code, id) {
		return
	}
	if err := p.PostRepository.UpdateById(ctx, id, gormplus.Update().WithColumns(
		dao.SysPostEntity.Name.Value(req.Name),
		dao.SysPostEntity.Code.Value(req.Code),
		dao.SysPostEntity.Sort.Value(req.Sort),
		dao.SysPostEntity.Status.Value(req.Status),
		dao.SysPostEntity.Remark.Value(req.Remark),
	).Build()); err != nil {
		p.Fail(ctx, err, "修改岗位失败")
		return
	}
	p.Success(ctx, "修改成功")
}

// GetPostPageApi 分页获取岗位
// @Summary 分页获取岗位
// @Description 根据岗位名称、编码、状态分页查询岗位列表
// @Tags 岗位中心
// @Accept json
// @Produce json
// @Param data body dto.SysPostPageDTO true "岗位分页查询参数"
// @Success 200 {array} vo.SysPostVO "岗位分页列表"
// @Router /api/v1/admin/post/page [post]
func (p Post) GetPostPageApi(ctx *gin.Context) {
	var req dto.SysPostPageDTO
	if !p.BindAndValidateJSON(ctx, &req) {
		return
	}
	list, total, err := p.PostRepository.FindPageByWrapper(ctx, req.PageNumber, req.PageSize, func(g gormplus.IGenWrapper[dao.ISysPostEntityDo]) {
		g.WhereIf(req.Name != "", dao.SysPostEntity.Name.Like("%"+req.Name+"%")).
			WhereIf(req.Code != "", dao.SysPostEntity.Code.Like("%"+req.Code+"%")).
			WhereIf(req.Status > 0, dao.SysPostEntity.Status.Eq(req.Status))
	}, gormplus.QueryOpt().Order(dao.SysPostEntity.Sort, dao.SysPostEntity.ID).Build())
	if err != nil {
		p.Fail(ctx, err, "获取岗位分页失败")
		return
	}
	p.BuildPageData(ctx, k.Map(list, func(item *model.SysPostEntity, index int) vo.SysPostVO {
		return p.toVO(item)
	}), total)
}

// GetPostListApi 获取岗位列表
// @Summary 获取岗位列表
// @Description 根据名称、编码和状态获取岗位列表
// @Tags 岗位中心
// @Accept json
// @Produce json
// @Param name query string false "岗位名称"
// @Param code query string false "岗位编码"
// @Param status query int false "状态：1正常，2禁用，0全部"
// @Success 200 {array} vo.SysPostVO "岗位列表"
// @Router /api/v1/admin/post/list [get]
func (p Post) GetPostListApi(ctx *gin.Context) {
	name := ctx.DefaultQuery("name", "")
	code := ctx.DefaultQuery("code", "")
	status := cast.ToInt64(ctx.DefaultQuery("status", "0"))
	list, err := p.PostRepository.FindListByWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysPostEntityDo]) {
		g.WhereIf(name != "", dao.SysPostEntity.Name.Like("%"+name+"%")).
			WhereIf(code != "", dao.SysPostEntity.Code.Like("%"+code+"%")).
			WhereIf(status > 0, dao.SysPostEntity.Status.Eq(status))
	}, gormplus.QueryOpt().Order(dao.SysPostEntity.Sort, dao.SysPostEntity.ID).Build())
	if err != nil {
		p.Fail(ctx, err, "获取岗位列表失败")
		return
	}
	p.Success(ctx, k.Map(list, func(item *model.SysPostEntity, index int) vo.SysPostVO {
		return p.toVO(item)
	}))
}

// GetPostDetailByIdApi 获取岗位详情
// @Summary 获取岗位详情
// @Description 根据岗位 ID 获取岗位详情，自动应用数据权限
// @Tags 岗位中心
// @Accept json
// @Produce json
// @Param id path int true "岗位ID"
// @Success 200 {object} vo.SysPostVO "岗位详情"
// @Router /api/v1/admin/post/detail/{id} [get]
func (p Post) GetPostDetailByIdApi(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	if id <= 0 {
		p.Fail(ctx, errors.New("参数id错误"), "参数id错误")
		return
	}
	entity, err := p.PostRepository.FindById(ctx, id)
	if err != nil {
		p.Fail(ctx, err, "岗位不存在或无权限访问")
		return
	}
	p.Success(ctx, p.toVO(entity))
}

func (p Post) checkCodeUnique(ctx *gin.Context, code string, excludeID int64) bool {
	conds := []gen.Condition{
		dao.SysPostEntity.Code.Eq(code),
	}
	if excludeID > 0 {
		conds = append(conds, dao.SysPostEntity.ID.Neq(excludeID))
	}
	exists, err := p.PostRepository.Exists(ctx, gormplus.QueryOpt().Where(conds...).Build())
	if err != nil {
		p.Fail(ctx, err, "岗位编码重复校验失败")
		return false
	}
	if exists {
		p.Fail(ctx, errors.New("岗位编码已存在"), "岗位编码已存在")
		return false
	}
	return true
}

func (p Post) toVO(entity *model.SysPostEntity) vo.SysPostVO {
	return p.PostMapper.EntityToVO(entity)
}

func NewPost(baseApi *base.BaseApi) IPost {
	return Post{
		BaseApi:               baseApi,
		PostRepository:        repository.NewSysPostRepository(),
		AccountPostRepository: repository.NewSysAccountPostRepository(),
		PostMapper:            mapper.NewSysPostMapper(),
	}
}
