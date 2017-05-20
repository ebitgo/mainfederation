package models

import (
	_webdb "github.com/ebitgo/ebitgo.com/models/databasemodels"
)

// LedgerDBControl Ledgercn数据库控制器
type LedgerDBControl struct {
	BaseManager
}

// CreateLedgerDbControl 创建Ledgercn数据库控制器
func CreateLedgerDbControl(conf DatabaseInfo, index int) IDataBase {
	ret := &LedgerDBControl{}
	ret.Indexof = index
	ret.SetDatabaseConfig(conf)
	ret.InitEngine("ledgercn")
	ret.InitOrm(new(_webdb.WalletDbT))
	ret.InitOperation(new(_webdb.WalletTableOperation))
	return ret
}
