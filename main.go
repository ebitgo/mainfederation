package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	_m "github.com/ebitgo/mainfederation/models"
	_ "github.com/ebitgo/mainfederation/routers"
)

func main() {
	mode := beego.AppConfig.String("runmode")
	_m.LoggerInstance = _m.NewLoggerInstance(fmt.Sprintf("log/federation.%s", time.Now().Format("2006-01-02_15.04.05.000")))
	_m.LoggerInstance.OpenDebug = strings.Compare("dev", mode) == 0
	_m.LoggerInstance.SetLogFunCallDepth(4)
	_m.LoggerInstance.Info("Start federation service....\r\n")
	conf := _m.DatabaseInfo{
		Host:     "localhost",
		Port:     "3306",
		UserName: "root",
		Password: "1234",
	}
	_m.DatabaseInstance = _m.CreateDBInstance(conf)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "DELETE", "PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
	}))
	beego.Run()
}
