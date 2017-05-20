package models

import (
	_wxdb "jojopoper/weixinSer/models/db"

	_webdb "github.com/ebitgo/ebitgo.com/models/databasemodels"
)

// DatabaseInstance 数据库访问唯一实例
var DatabaseInstance *DatabaseManager

const (
	LedgercnIndex = 0
	WechatIndex   = 1
)

// DatabaseManager 数据库控制器
type DatabaseManager struct {
	conf        DatabaseInfo
	controllers []IDataBase
}

// CreateDBInstance 创建数据库控制器唯一实例
func CreateDBInstance(conf DatabaseInfo) *DatabaseManager {
	ret := new(DatabaseManager)
	return ret.Init(conf)
}

// Init 初始化
func (ths *DatabaseManager) Init(conf DatabaseInfo) *DatabaseManager {
	ths.conf = conf
	ths.controllers = make([]IDataBase, 2)
	ths.controllers[LedgercnIndex] = CreateLedgerDbControl(conf, LedgercnIndex)
	ths.controllers[WechatIndex] = CreateWxDbControl(conf, WechatIndex)
	return ths
}

// GetWalletInfo ledgercn.com 中的用户查找
func (ths *DatabaseManager) GetWalletInfo(addr, nickname string) ([]*_webdb.WalletDbT, error) {
	ret := make(map[string][]*_webdb.WalletDbT)
	opera := ths.controllers[LedgercnIndex].GetOperation(_webdb.DB_WALLET_INFO_OPERATION)
	err := opera.Quary(_webdb.QT_SEARCH_RECORD, ret, addr, nickname)
	if err == nil {
		walletinfo := ret["RESULT"]
		return walletinfo, nil
	}
	return nil, err
}

// GetUserInfo 在微信数据库中查找
func (ths *DatabaseManager) GetUserInfo(addr, nickname string) (*_wxdb.TUserInfo, error) {
	opera := ths.controllers[WechatIndex].GetOperation(_wxdb.DbUserInfoOperation)
	ui := &_wxdb.TUserInfo{
		UniqueId:   nickname,
		PublicAddr: addr,
	}
	err := opera.Quary(_wxdb.QtGetRecord, ui)
	return ui, err
}
