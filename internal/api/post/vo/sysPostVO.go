package vo

// SysPostVO 岗位视图对象
type SysPostVO struct {
	ID        int64  `json:"id"`        // 主键id
	Name      string `json:"name"`      // 岗位名称
	Code      string `json:"code"`      // 岗位编码，全局唯一
	Sort      int64  `json:"sort"`      // 排序，值越小越靠前
	Status    int64  `json:"status"`    // 状态：1正常，2禁用
	Remark    string `json:"remark"`    // 备注
	CreatedAt int64  `json:"createdAt"` // 创建时间，Unix秒
	UpdatedAt int64  `json:"updatedAt"` // 更新时间，Unix秒
	CreatedBy int64  `json:"createdBy"` // 创建人账号id
	UpdatedBy int64  `json:"updatedBy"` // 更新人账号id
}
