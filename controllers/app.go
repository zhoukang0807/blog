// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"myblog/models"
	"regexp"
	"strings"
)

var langTypes []string                                                              // Languages that are supported.
var digitsRegexp = regexp.MustCompile(`.*?<blockquote><p>(.*?)<\/p><\/blockquote>`) //提取文本中标签的内容

func init() {
	// Initialize language type list.
	langTypes = strings.Split(beego.AppConfig.String("lang_types"), "|")

	// Load locale files according to language types.
	for _, lang := range langTypes {
		beego.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file:", err)
			return
		}
	}
}

// baseController represents base router for all other app routers.
// It implemented some methods for the same implementation;
// thus, it will be embedded into other routers.
type baseController struct {
	beego.Controller // Embed struct that has stub implementation of the interface.
	i18n.Locale      // For i18n usage when process data and render template.
}
type AppController struct {
	baseController
}

// Prepare implemented Prepare() method for baseController.
// It's used for language option check and setting.
func (this *baseController) Prepare() {
	// Reset language option.
	this.Lang = "" // This field is from i18n.Locale.

	// 1. Get language information from 'Accept-Language'.
	al := this.Ctx.Request.Header.Get("Accept-Language")
	if len(al) > 4 {
		al = al[:5] // Only compare first 5 letters.
		if i18n.IsExist(al) {
			this.Lang = al
		}
	}

	// 2. Default language is English.
	if len(this.Lang) == 0 {
		this.Lang = "en-US"
	}

	// Set template level language option.
	this.Data["Lang"] = this.Lang
}

// Get implemented Get() method for LoginController.
func (this *AppController) Get() {
	this.TplName = "index.html"
	this.Data["IsHome"] = true
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	_, topics, err := models.GetAllTopics("", this.GetString("label"), this.GetString("cate"), "created", "-", 5, 0)

	for _, topic := range topics {
		topic.Content = getMainContent(topic.Content)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Data["Topics"] = topics
	stickTopics, err := models.GetAllTopicsByStick()
	if err != nil {
		beego.Error(err)
	}
	this.Data["StickTopics"] = stickTopics

	categroies, err := models.GetAllCategory()
	if err != nil {
		beego.Error(err)
	}
	this.Data["Categories"] = categroies

}

func (this *AppController) Game() {
	this.TplName = "game.html"
	this.Data["IsGame"] = true
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}

func (this *AppController) MessageBoard() {
	this.TplName = "messageboard.html"
	this.Data["IsBoard"] = true
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}

func getMainContent(topic string) string {
	cont := digitsRegexp.FindStringSubmatch(topic) //获取截取后的slience
	var conts []string
	if len(cont) > 0 {
		conts = cont[1:]
	} else {
		return ""
	}

	var strcontent string
	for _, topic := range conts {
		strcontent += topic
	}
	return strcontent

}
