package controllers

import (
	"github.com/Unknwon/com" //通用方法
	"github.com/astaxie/beego"
	"myblog/models"
	"os" //文件操作
	"path"
	"strconv"
	"strings"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type TopicController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for LoginController.
func (this *TopicController) Get() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.TplName = "topic.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	this.Data["IsTopic"] = true
	_, topics, err := models.GetAllTopics("", "", "", "", "", 5, 0)
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Topics"] = topics
	}

}
func (this *TopicController) Post() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	//解析表单
	title := this.GetString("title")
	content := this.GetString("content")
	beego.Info(content)
	tid := this.GetString("tid")
	category := this.GetString("category")
	ctype := this.GetString("type")
	label := this.GetString("label")
	//获取附件
	f, fh, err := this.GetFile("attachment")
	var attachment string
	if err != nil {
		beego.Error(err)
		attachment = ""
	} else {
		if fh != nil {
			//保存附件
			attachment = fh.Filename
			beego.Info(attachment)
			err = this.SaveToFile("attachment", path.Join("attachment/topic", attachment))
			if err != nil {
				beego.Error(err)
			}
		}
	}
	defer f.Close()
	//----------------图片压缩
	var savePath string = "attachment/topic"
	width := 200
	var themePath string = "attachment/topic/" + strconv.Itoa(width) + "w"
	if !com.IsExist(themePath) {
		os.MkdirAll(themePath, os.ModePerm)
	}
	saveThemeFile(attachment, savePath, themePath, width)
	//----------压缩完成
	if len(tid) == 0 {
		err = models.AddTopic(ctype, title, category, label, content, attachment)
	} else {
		err = models.ModifyTopic(tid, ctype, title, category, label, content, attachment)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/article", 302)

}
func (this *TopicController) Add() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.TplName = "topic_add.html"
}
func (this *TopicController) Category() {
	ctype, err := strconv.ParseInt(this.GetString("type"), 10, 64) //id 十进制 int64
	if err != nil {
		beego.Error(err)
	}
	categorys, err := models.GetAllCategoryByType(ctype)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"success": false, "message": "faild", "data": nil}
	} else {
		this.Data["json"] = map[string]interface{}{"success": true, "message": "success", "data": categorys}
	}

	this.ServeJSON()
}

func (this *TopicController) Modify() {
	this.TplName = "topic_modify.html"
	tid := this.GetString("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	this.Data["Topic"] = topic
}

func (this *TopicController) Delete() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	err := models.DeleteTopic(this.GetString("tid"))
	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/article", 302)
}
func (this *TopicController) Stick() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	tid := this.GetString("tid")
	stick, err := models.StickTopic(tid)

	if err != nil {
		this.Data["json"] = map[string]interface{}{"success": false, "message": "faild"}
	} else {
		this.Data["json"] = map[string]interface{}{"success": true, "message": "success", "stick": stick}
	}

	if err != nil {
		beego.Error(err)
	}

	this.ServeJSON()
}
func (this *TopicController) View() {
	this.TplName = "topic_view.html"
	//reqUrl := this.Ctx.Request.RequestURI 获取前一个页面的url
	// i := strings.LastIndex(reqUrl, "/")
	// tis := reqUrl[i+1:]
	topic, err := models.GetTopic(this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	metaContent := getMainContent(topic.Content)
	this.Data["MetaContent"] = metaContent
	this.Data["Topic"] = topic
	this.Data["Labels"] = strings.Split(topic.Labels, " ")
	_, topics, err := models.GetAllTopics("", "", "", "views", "-", 6, 0)
	if err != nil {
		beego.Error(err)
	}
	beego.Error(topic.Ctype)
	this.Data["HotTopics"] = topics
	// replies, err := models.GetAllReplies(this.Ctx.Input.Param("0"))
	// if err != nil {
	// 	beego.Error(err)
	// 	return
	// }
	// this.Data["Replies"] = replies
	RandTopics, err := models.GetAllTopicsRand()
	if err != nil {
		beego.Error(err)
	}

	this.Data["RandTopics"] = RandTopics
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	this.Data["IsArticle"] = true
}
