package initialize

import (
	"gin_admin_api/global"
	"gin_admin_api/model"
)

func InitDataSource() {
	global.DB.AutoMigrate(&model.AccountEntity{})
}
