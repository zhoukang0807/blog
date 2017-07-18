package controllers

import (
	"github.com/astaxie/beego"
	"math"
	"myblog/models"
	"strconv"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type ArticleController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

func (this *ArticleController) Get() {
	this.TplName = "article.html"

	pagesize := int64(10)
	cnt, topics, err := models.GetAllTopics("", "", "", "created", "-", pagesize, 0)
	for _, topic := range topics {
		topic.Content = getMainContent(topic.Content)
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
	if err != nil {
		beego.Error(err)
	}
	this.Data["IsArticle"] = true
	this.Data["Topics"] = topics
	this.Data["Count"] = cnt
	this.Data["PageArray"] = array1
	this.Data["Next"] = 2
	this.Data["CurrentPage"] = 1
	this.Data["TotalPage"] = totalpage
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	//热门博客
	_, Hottopics, err := models.GetAllTopics("", "", "", "views", "-", 6, 0)
	if err != nil {
		beego.Error(err)
	}
	this.Data["HotTopics"] = Hottopics
}

func (this *ArticleController) Aplit() {
	this.TplName = "article.html"
	searchText := this.GetString("q")
	beego.Info("查询条件为" + searchText)
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

	cnt, topics, err := models.GetAllTopics(searchText, this.GetString("label"), this.GetString("cate"), "created", "-", pagesize, start)
	for _, topic := range topics {
		topic.Content = getMainContent(topic.Content)
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

	this.Data["IsArticle"] = true
	this.Data["Topics"] = topics
	this.Data["Count"] = cnt
	this.Data["PageArray"] = array1
	this.Data["CurrentPage"] = startpage
	this.Data["Next"] = startpage + 1
	this.Data["Prev"] = startpage - 1
	this.Data["SearchTxt"] = searchText
	this.Data["TotalPage"] = totalpage
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	//热门博客
	_, Hottopics, err := models.GetAllTopics("", "", "", "views", "-", 6, 0)
	if err != nil {
		beego.Error(err)
	}
	this.Data["HotTopics"] = Hottopics
}
func (this *ArticleController) Search() {
	SearchText := this.GetString("search")
	this.Redirect("/article/aplit/1?q="+SearchText, 302)
}
