package routers

import (
	"github.com/astaxie/beego"
	"github.com/ebitgo/mainfederation/controllers"
)

func init() {
	beego.Router("/federation", &controllers.FederationController{})
}
