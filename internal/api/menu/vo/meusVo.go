package vo

type MenusVO struct {
	ID            int64  `json:"id"`            // 主键id
	Title         string `json:"title"`         // 名称:按钮标题,或菜单标题
	URL           string `json:"url"`           // 按钮请求url,或菜单路由
	Icon          string `json:"icon"`          // 菜单小图标
	ResourcesType int64  `json:"resourcesType"` // 类型:1表示目录,2表示菜单,3表示接口
	IsCache       int64  `json:"isCache"`       // 是否缓存:1表示缓存:2不缓存
	IsHidden      int64  `json:"isHidden"`      // 是否隐藏:1表示不隐藏,2表示隐藏
	IsLink        int64  `json:"isLink"`        // 是否为外部链接:1表示不是,2表示是
	ParentID      int64  `json:"parentId"`      // 上一级id
	Sort          int64  `json:"sort"`          // 菜单,或按钮排序
}
