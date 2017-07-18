package controllers

import (
	"github.com/Unknwon/com" //通用方法
	"github.com/astaxie/beego"
	"graphics-go/graphics"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"myblog/models"
	"os" //文件操作
	"path"
	"strconv"
	"strings"
)

// LoginController handles the welcome screen that allows user to pick a technology and username.
type GalleryController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

func (this *GalleryController) Delete() {
	id := this.GetString("id")
	beego.Info("图片ID:" + id)
	err := models.DelGallery(id)
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/gallery", 302)
}

func (this *GalleryController) Add() {
	this.TplName = "gallery_add.html"
	this.Data["IsGallerys_Add"] = true
}

func (this *GalleryController) Get() {
	this.TplName = "gallery.html"
	gallerys, err := models.GetAllGallery()
	if err != nil {
		beego.Error(err)
	}
	this.Data["IsGallery"] = true
	this.Data["Gallerys"] = gallerys
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}

func (this *GalleryController) Post() {
	ck, err := this.Ctx.Request.Cookie("uname")
	if err != nil {
		return
	}
	uname := ck.Value
	beego.Info("上传文件服务开始----")
	var savePath string = "attachment/" + uname + "/gallery"
	if !com.IsExist(savePath) {
		os.MkdirAll(savePath, os.ModePerm)
	}
	f1, fh, err := this.GetFile("webfiles")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		//保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = this.SaveToFile("webfiles", path.Join(savePath, attachment))
		if err != nil {
			beego.Error(err)
		}
	}
	defer f1.Close()
	beego.Info("上传文件服务结束----")

	width := 300
	var themePath string = "attachment/" + uname + "/gallery/" + strconv.Itoa(width) + "w"
	if !com.IsExist(themePath) {
		os.MkdirAll(themePath, os.ModePerm)
	}
	saveThemeFile(attachment, savePath, themePath, width)
	err = models.AddGallery(attachment, uname)
	if err != nil {
		beego.Error(err)
	}
	beego.Info("上传" + attachment + "成功")
	this.Redirect("/gallery", 302)
}

/**
 * 压缩并保存到目录
 */
func saveThemeFile(attachment, savePath, themePath string, width int) {
	beego.Info("压缩文件服务开始----")
	pic1, format, name := isPictureFormat(attachment)
	beego.Info(pic1 + "--" + format + "---" + name)
	if strings.ToLower(format) == "jpg" || strings.ToLower(format) == "jpeg" {
		f3, err := os.Create(path.Join(themePath, attachment))
		if err != nil {
			beego.Error(err)
		}
		defer f3.Close()
		img := ImageJPEG(path.Join(savePath, attachment))
		if img == nil {
			beego.Error("图片未找到")
		}
		bound := img.Bounds()
		dx := bound.Dx()
		dy := bound.Dy()

		bounds := image.Rect(0, 0, width, width*dy/dx)
		m := image.NewRGBA(bounds)
		err = graphics.Scale(m, img)
		if err != nil {
			beego.Error(err)
		}

		err = jpeg.Encode(f3, m, &jpeg.Options{90})
		if err != nil {
			beego.Error(err)
		}

	} else if strings.ToLower(format) == "png" {
		f3, err := os.Create(path.Join(themePath, attachment))
		if err != nil {
			beego.Error(err)
		}
		defer f3.Close()
		img := ImagePNG(path.Join(savePath, attachment))
		if img == nil {
			beego.Error("图片未找到")
		}
		bound := img.Bounds()
		dx := bound.Dx()
		dy := bound.Dy()

		bounds := image.Rect(0, 0, width, width*dy/dx)
		m := image.NewRGBA(bounds)
		err = graphics.Scale(m, img)
		if err != nil {
			beego.Error(err)
		}

		err = png.Encode(f3, m)
		if err != nil {
			beego.Error(err)
		}
	} else if strings.ToLower(format) == "gif" {
		f3, err := os.Create(path.Join(themePath, attachment))
		if err != nil {
			beego.Error(err)
		}
		defer f3.Close()
		img := ImageGIF(path.Join(savePath, attachment))
		if img == nil {
			beego.Error("图片未找到")
		}
		bound := img.Bounds()
		dx := bound.Dx()
		dy := bound.Dy()

		bounds := image.Rect(0, 0, width, width*dy/dx)
		m := image.NewRGBA(bounds)
		err = graphics.Scale(m, img)
		if err != nil {
			beego.Error(err)
		}

		err = gif.Encode(f3, m, nil)
		if err != nil {
			beego.Error(err)
		}
	}
	beego.Info("压缩文件结束----")
}

/** 是否是图片 */
func isPictureFormat(path string) (string, string, string) {
	temp := strings.Split(path, ".")
	if len(temp) <= 1 {
		return "", "", ""
	}
	mapRule := make(map[string]int64)
	mapRule["jpg"] = 1
	mapRule["png"] = 1
	mapRule["jpeg"] = 1
	mapRule["gif"] = 1
	mapRule["bmp"] = 1
	// fmt.Println(temp[1]+"---")
	/** 添加其他格式 */
	if mapRule[strings.ToLower(temp[1])] == 1 {
		println(temp[1])
		return path, temp[1], temp[0]
	} else {
		return "", "", ""
	}
}

// 读取JPEG图片返回image.Image对象
func ImageJPEG(ph string) image.Image {
	// 打开图片文件
	f, fileErr := os.Open(ph)
	if fileErr != nil {
		return nil
	}
	// 退出时关闭文件
	defer f.Close()
	// 解码
	j, jErr := jpeg.Decode(f)
	if jErr != nil {
		return nil
	}
	// 返回解码后的图片
	return j
}

// 读取PNG图片返回image.Image对象
func ImagePNG(ph string) image.Image {
	// 打开图片文件
	f, fileErr := os.Open(ph)
	if fileErr != nil {
		return nil
	}
	// 退出时关闭文件
	defer f.Close()
	// 解码
	p, pErr := png.Decode(f)
	if pErr != nil {
		return nil
	}
	// 返回解码后的图片
	return p
}

// 读取PNG图片返回image.Image对象
func ImageGIF(ph string) image.Image {
	// 打开图片文件
	f, fileErr := os.Open(ph)
	if fileErr != nil {
		return nil
	}
	// 退出时关闭文件
	defer f.Close()
	// 解码
	p, pErr := gif.Decode(f)
	if pErr != nil {
		return nil
	}
	// 返回解码后的图片
	return p
}
