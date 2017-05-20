package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	_m "github.com/ebitgo/mainfederation/models"
)

//FederationController 联邦服务器
type FederationController struct {
	beego.Controller
}

// Get get function
func (ths *FederationController) Get() {
	ths.setCtxHeader()
	msg := &_m.ResponseMsg{}
	msg.DecodeParam(ths.Input())
	if msg.Status == 200 {
		msg.Execute()
	}

	if msg.Status == 200 {
		ths.Data["json"] = msg.Msg
		_m.LoggerInstance.InfoPrint("[%s] [%+v] \r\nget message success!\r\n", ths.Ctx.Input.IP(), ths.Input())
	} else {
		ths.Data["json"] = map[string]interface{}{
			"detail": msg.ErrMsg,
		}
		http.Error(ths.Ctx.ResponseWriter, "", msg.Status)
		_m.LoggerInstance.InfoPrint("[%s] [%+v] \r\nget message failure! [errmsg : %s]\r\n", ths.Ctx.Input.IP(), ths.Input(), msg.ErrMsg)
	}
	ths.ServeJSON(false)
}

func (ths *FederationController) setCtxHeader() {
	// ths.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Accept")
	// ths.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	// ths.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Methods", "GET")
}
