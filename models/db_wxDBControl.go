package models

import (
	_wxdb "jojopoper/weixinSer/models/db"
)

// WxDBControl 微信数据库控制器
type WxDBControl struct {
	BaseManager
}

// CreateWxDbControl 创建微信数据库控制器
func CreateWxDbControl(conf DatabaseInfo, index int) IDataBase {
	ret := &WxDBControl{}
	ret.Indexof = index
	ret.SetDatabaseConfig(conf)
	ret.InitEngine("weixin")
	ret.InitOrm(new(_wxdb.TUserInfo))
	ret.InitOperation(new(_wxdb.UserInfoOperation))
	return ret
}
