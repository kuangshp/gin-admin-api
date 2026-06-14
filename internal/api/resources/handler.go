package resources

import (
	"gin-admin-api/internal/api/base"
	"github.com/gin-gonic/gin"
)

type IResources interface {
	CreateResourcesApi(ctx *gin.Context)      // 创建资源
	DeleteResourcesByIdApi(ctx *gin.Context)  // 根据主键删除资源
	ModifyResourcesByIdApi(ctx *gin.Context)  // 根据主键修改资源
	GetResourcesTreePageApi(ctx *gin.Context) // 获取资源树
	GetResourcesCatalogApi(ctx *gin.Context)  //  查询目录或者目录和菜单
	GetResourcesListApi(ctx *gin.Context)     // 获取全部的目录及菜单
}

type Resources struct {
	*base.BaseApi
}

func (r Resources) CreateResourcesApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r Resources) DeleteResourcesByIdApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r Resources) ModifyResourcesByIdApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r Resources) GetResourcesTreePageApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r Resources) GetResourcesCatalogApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r Resources) GetResourcesListApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewResources(baseApi *base.BaseApi) IResources {
	return Resources{
		BaseApi: baseApi,
	}
}
