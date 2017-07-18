package controllers

import (
	"myblog/models"

	"github.com/astaxie/beego"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type CategoryController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for LoginController.
func (this *CategoryController) Get() {

	this.TplName = "category.html"
	this.Data["IsCategory"] = true
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	var err error
	this.Data["Categories"], err = models.GetAllCategory() //this.Data是定义好的对象无法重新定义 不能用：
	if err != nil {
		beego.Error(err)
	}
}

func (this *CategoryController) Add() {
	name := this.GetString("name")
	ctype := this.GetString("type")
	if len(name) == 0 {
		return
	}
	err := models.AddCategory(name, ctype)
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/category", 301)
}
func (this *CategoryController) Del() {
	id := this.GetString("id")
	if len(id) == 0 {
		return
	}
	err := models.DelCategory(id)
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/category", 301)
	return
}
