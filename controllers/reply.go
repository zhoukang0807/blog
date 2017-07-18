package controllers

import (
	"github.com/astaxie/beego"
	"myblog/models"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type ReplyController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

func (r *ReplyController) Add() {
	tid := r.GetString("tid")
	err := models.AddReply(tid, r.GetString("nickname"), r.GetString("content"))
	if err != nil {
		beego.Error(err)
	}

	r.Redirect("/topic/view/"+tid, 302)
}
func (r *ReplyController) Delete() {
	if !checkAccount(r.Ctx) {
		return
	}
	tid := r.GetString("tid")
	rid := r.GetString("rid")
	err := models.DeleteReply(rid)
	if err != nil {
		beego.Error(err)
	}

	r.Redirect("/topic/view/"+tid, 302)
}
