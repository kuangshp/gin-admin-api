package datascope

type Context struct {
	AccountID int64
	IsAdmin   int64
	DeptID    int64
	RoleIDs   []int64
	DataScope int64
	DeptIDs   []int64
}

func (c Context) IsSuperAdmin() bool {
	return c.IsAdmin == 1
}
