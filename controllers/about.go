package controllers

import (
	"github.com/Unknwon/com" //通用方法
	"github.com/astaxie/beego"
	"myblog/models"
	"os"
	"path"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type AboutController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for LoginController.
func (this *AboutController) Get() {
	this.TplName = "about.html"
	ck, err := this.Ctx.Request.Cookie("uname")
	var uname string
	if err != nil {
		beego.Error(err)
		uname = "admin"
	} else {
		uname = ck.Value
	}

	about, err := models.GetAbout(uname)
	this.Data["IsAbout"] = true
	//this.Data["IsNull"] = about == nil
	this.Data["About"] = about
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}
func (this *AboutController) Add() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.TplName = "about_add.html"
	this.Data["IsAbout"] = true
}
func (this *AboutController) Modify() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	ck, err := this.Ctx.Request.Cookie("uname")
	var uname string
	if err != nil {
		beego.Error(err)
		uname = "admin"
	} else {
		uname = ck.Value
	}
	about, err := models.GetAbout(uname)
	this.Data["About"] = about
	this.TplName = "about_modify.html"
	this.Data["IsAbout"] = true
}
func (this *AboutController) Post() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	//解析表单
	athor := this.GetString("athorName")
	tid := this.GetString("tid")
	content := this.GetString("content")
	//获取附件
	f, fh, err := this.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	ck, err := this.Ctx.Request.Cookie("uname")
	if err != nil {
		return
	}
	uname := ck.Value
	beego.Info("上传文件服务开始----")
	var savePath string = "attachment/" + uname + "/about"
	if !com.IsExist(savePath) {
		os.MkdirAll(savePath, os.ModePerm)
	}
	var attachment string
	if fh != nil {
		//保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = this.SaveToFile("attachment", path.Join(savePath, attachment))
		if err != nil {
			beego.Error(err)
		}
	}
	if len(tid) == 0 {
		err = models.AddAbout(uname, athor, attachment, content)

	} else {
		err = models.ModifyAbout(tid, uname, athor, attachment, content)
	}
	if err != nil {
		beego.Error(err)
	}
	beego.Info("上传文件服务结束----")
	defer f.Close()
	this.Redirect("/about", 302)
}
