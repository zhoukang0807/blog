package controllers

import (
	"github.com/astaxie/beego"
	"math"
	"myblog/models"
	"strconv"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type SpeackController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

func (this *SpeackController) Get() {
	this.TplName = "speack.html"
	pagesize := int64(10)
	cnt, speacks, err := models.GetAllSpeack(pagesize, 0)
	if err != nil {
		beego.Error(err)
	}
	totalpage := int(math.Ceil(float64(cnt) / float64(pagesize)))
	var array1 = [5]int{}
	if totalpage < 5 {
		for i := 0; i < totalpage; i++ {
			array1[i] = i + 1
		}
	} else {
		for i := 0; i < 5; i++ {
			array1[i] = i + 1
		}
	}

	this.Data["Speacks"] = speacks
	this.Data["IsSpeack"] = true
	this.Data["Count"] = cnt
	this.Data["PageArray"] = array1
	this.Data["Next"] = 2
	this.Data["CurrentPage"] = 1
	this.Data["TotalPage"] = totalpage
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}
func (this *SpeackController) Split() {
	this.TplName = "speack.html"
	pagesize := int64(10)
	startpage, err := strconv.ParseInt(this.Ctx.Input.Param("0"), 10, 64)
	start := startpage
	if err != nil {
		beego.Error(err)
		start = 0
	} else {
		start = (start - 1) * pagesize
		beego.Error(start)
	}

	cnt, speacks, err := models.GetAllSpeack(pagesize, start)
	if err != nil {
		beego.Error(err)
	}
	totalpage := int(math.Ceil(float64(cnt) / float64(pagesize)))

	var array1 = [5]int{}
	chazhi := totalpage - int(startpage)
	if totalpage < 5 {
		for i := 0; i < totalpage; i++ {
			array1[i] = i + 1
		}
	} else {
		if chazhi >= 5 {
			for i := 0; i < 5; i++ {
				array1[i] = int(startpage) + i
			}
		} else {
			for i := 4; i >= 0; i-- {
				array1[i] = totalpage - (4 - i)
			}
		}
	}

	this.Data["Speacks"] = speacks
	this.Data["IsSpeack"] = true
	this.Data["Count"] = cnt
	this.Data["PageArray"] = array1
	this.Data["CurrentPage"] = startpage
	this.Data["Next"] = startpage + 1
	this.Data["Prev"] = startpage - 1
	this.Data["TotalPage"] = totalpage
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}

func (this *SpeackController) Post() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	ck, err := this.Ctx.Request.Cookie("uname")
	if err != nil {
		return
	}
	uname := ck.Value
	content := this.GetString("content")
	err = models.AddSpeack(uname, content)
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/speack", 302)
}
