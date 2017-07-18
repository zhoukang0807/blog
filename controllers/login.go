package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type LoginController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for LoginController.
func (this *LoginController) Get() {
	isExit := this.GetString("exit") == "true"
	if isExit {
		this.Ctx.SetCookie("uname", "", -1, "/")
		this.Ctx.SetCookie("passwd", "", -1, "/")
		this.Redirect("/", 301)
		return
	}
	this.TplName = "login.html"
	this.Data["IsLogin"] = true
}

func (this *LoginController) Post() {
	uname := this.GetString("uname")
	passwd := this.GetString("passwd")
	autoLogin := this.GetString("autoLogin") == "on"
	if beego.AppConfig.String("uname") == uname && beego.AppConfig.String("passwd") == passwd {
		maxAge := 0
		if autoLogin {
			maxAge = 1<<31 - 1
		}
		this.Ctx.SetCookie("uname", uname, maxAge, "/")
		this.Ctx.SetCookie("passwd", passwd, maxAge, "/")
	}

	this.Redirect("/", 301)
	return
}

func checkAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}
	uname := ck.Value
	ck, err = ctx.Request.Cookie("passwd") //获取请求的cookie
	if err != nil {
		return false
	}
	passwd := ck.Value
	return beego.AppConfig.String("uname") == uname && beego.AppConfig.String("passwd") == passwd
}
