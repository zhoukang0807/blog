package main

import (
	"myblog/controllers"
	"myblog/models"

	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
)

func init() {
	models.RegisterDB()
}

const (
	APP_VER = "0.1.7.0113"
)

func main() {
	orm.Debug = true
	orm.RunSyncdb("default", false, true) //1.数据库 2.是否删掉重建 3.打印
	beego.Info(beego.BConfig.AppName, APP_VER)
	// 注册相关路由
	beego.Router("/", &controllers.AppController{})
	beego.Router("/game", &controllers.AppController{}, "get:Game")

	beego.Router("/messageboard", &controllers.AppController{}, "get:MessageBoard")
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.AutoRouter(&controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.Router("/article", &controllers.ArticleController{})
	beego.Router("/gallery", &controllers.GalleryController{})
	beego.Router("/speack", &controllers.SpeackController{})
	beego.Router("/about", &controllers.AboutController{})
	beego.AutoRouter(&controllers.AboutController{})
	beego.AutoRouter(&controllers.SpeackController{})
	beego.AutoRouter(&controllers.GalleryController{})
	beego.AutoRouter(&controllers.ArticleController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.AutoRouter(&controllers.ReplyController{})
	//beego.Router("/reply", &controllers.ReplyController{})
	//beego.Router("/reply/add", &controllers.ReplyController{}, "post:Add")
	//beego.Router("/reply/delete", &controllers.ReplyController{}, "")
	//beego.Router("/reg", &controllers.RegisterController{}, "get:Get;post:Reg")

	//创建附件目录
	os.Mkdir("attachment", os.ModePerm)
	//附件作为一个控制器来处理
	beego.Router("/attachment/:all", &controllers.AttachController{})
	beego.Router("/attachment/:all/:all", &controllers.AttachController{})
	//作为静态文件处理附件
	beego.SetStaticPath("/attachment", "attachment")
	// 注册国际化
	beego.AddFuncMap("i18n", i18n.Tr)

	beego.Run()

	//交叉编译
	// set GOPACH=amd64
	// set GOOS=linux
	// bee pack 打成的包扔到linux上解压给权限后就可以运行了
	//supervisorctl restart myblog 重启服务

}
