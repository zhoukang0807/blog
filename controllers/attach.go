package controllers

import (
	"github.com/astaxie/beego"
	"io"
	"net/url"
	"os"
	"path"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type AttachController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

func (this *AttachController) Get() {
	filePath, err := url.QueryUnescape(this.Ctx.Request.RequestURI[1:]) //url解码
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}
	defer f.Close()

	_, err = io.Copy(this.Ctx.ResponseWriter, f)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}
}

func (this *AttachController) Post() {
	beego.Error("富文本框上传文件开始")
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	ck, err := this.Ctx.Request.Cookie("uname")
	if err != nil {
		return
	}
	uname := ck.Value
	//获取附件
	_, fh, err := this.GetFile("editFileName")
	//创建附件目录
	os.Mkdir("attachment/"+uname, os.ModePerm)
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		//保存附件
		attachment = fh.Filename
		beego.Info(attachment)

		err = this.SaveToFile("editFileName", path.Join("attachment/"+uname, attachment))
		if err != nil {
			beego.Error(err)
		}
	}
	this.Ctx.WriteString("/attachment/" + uname + "/" + attachment)
	beego.Error("富文本框上传文件结束")
}
