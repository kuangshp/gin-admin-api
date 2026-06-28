package enum

const (
	DataScopeAll          int64 = 1 // 全部数据
	DataScopeDept         int64 = 2 // 本部门数据
	DataScopeDeptAndChild int64 = 3 // 本部门及下级数据
	DataScopeSelf         int64 = 4 // 仅本人数据
	DataScopeCustomDept   int64 = 5 // 自定义部门数据
)
