package models

import (
	"fmt"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	MySqlDriver    = "mysql"
	SqLiteDirver   = "sqlite3"
	PostgresDriver = "postgres"
)

// DatabaseInfo 数据库基本信息定义
type DatabaseInfo struct {
	Host     string
	Port     string
	UserName string
	Password string
}

// OperationInterface 数据库接口定义
type OperationInterface interface {
	Init(e *xorm.Engine) error
	GetKey() string
	Quary(qtype int, v ...interface{}) error
}

// IDataBase 数据库基本定义接口
type IDataBase interface {
	SetDatabaseConfig(dbConfig DatabaseInfo)
	InitEngine(alisaName string)
	InitOrm(table interface{})
	InitOperation(opera OperationInterface)
	GetOperation(key string) OperationInterface
}

// BaseManager 数据库管理实例定义
type BaseManager struct {
	Indexof     int
	dbConfig    DatabaseInfo
	DbEngine    *xorm.Engine
	Operations  map[string]OperationInterface
	dbType      string
	dataSrcName string
}

// SetDatabaseConfig 初始化数据库参数
func (ths *BaseManager) SetDatabaseConfig(dbConfig DatabaseInfo) {
	LoggerInstance.InfoPrint("[BaseManager:SetDatabaseConfig] set database config \r\n")
	ths.dbConfig = dbConfig
}

// InitEngine 初始化数据库引擎
func (ths *BaseManager) InitEngine(alisaName string) {
	// 读取配置文件中的数据库类型
	ths.dbType = beego.AppConfig.String("dbtype")

	ths.DbEngine = nil
	switch ths.dbType {
	case MySqlDriver:
		ths.dataSrcName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", //Asia%2FShanghai
			ths.dbConfig.UserName, ths.dbConfig.Password, ths.dbConfig.Host, ths.dbConfig.Port, alisaName)
	case SqLiteDirver:
		ths.dataSrcName = fmt.Sprintf("./%s.db?charset=utf8&loc=Asia/Shanghai", alisaName)
	case PostgresDriver:
		ths.dataSrcName = fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=verify-full",
			alisaName, ths.dbConfig.UserName, ths.dbConfig.Password, ths.dbConfig.Host, ths.dbConfig.Port)
	default:
		LoggerInstance.ErrorPrint("[BaseManager:InitEngine] Undefined db type = %s\r\n", ths.dbType)
		panic(1)
	}

	ths.DbEngine = ths.getEngine(ths.dbType)
	if ths.DbEngine == nil {
		return
	}

	isDebug := beego.AppConfig.String("runmode")
	if isDebug == "dev" {
		ths.DbEngine.ShowDebug = true
		ths.DbEngine.ShowInfo = true
		ths.DbEngine.ShowSQL = true
	}
	ths.DbEngine.ShowErr = true
	ths.DbEngine.ShowWarn = true
}

// InitOrm 初始化数据库表
func (ths *BaseManager) InitOrm(table interface{}) {
	err := ths.DbEngine.Sync(table)
	if err != nil {
		LoggerInstance.InfoPrint("[BaseManager:InitOrm] XORM Engine Sync is err %v\r\n", err)
		panic(1)
	}
}

// InitOperation 初始化操作
func (ths *BaseManager) InitOperation(opera OperationInterface) {
	if ths.Operations == nil {
		ths.Operations = make(map[string]OperationInterface)
	}

	opera.Init(ths.DbEngine)
	ths.Operations[opera.GetKey()] = opera
}

// GetOperation 读取Operation
func (ths *BaseManager) GetOperation(key string) (ret OperationInterface) {
	ret,_ = ths.Operations[key]
	return
}

func (ths *BaseManager) getEngine(dbDrv string) *xorm.Engine {
	ret, err := xorm.NewEngine(dbDrv, ths.dataSrcName)
	if err != nil {
		LoggerInstance.ErrorPrint("[BaseManager:getMySqlEngine] Create %s has error! \r\n\t%v\r\n", dbDrv, err)
		return nil
	}
	err = ret.Ping()
	if err != nil {
		LoggerInstance.ErrorPrint("[BaseManager:getMySqlEngine] Create %s Ping error! \r\n\t %v\r\n", dbDrv, err)
		return nil
	}
	return ret
}
