package models

import (
	//"github.com/Unknwon/com"       //通用方法
	"github.com/astaxie/beego/orm" //数据库操作
	//_ "github.com/mattn/go-sqlite3" //sqlite3驱动包
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql" //mysql驱动
	"os"                               //文件操作
	"path"                             //路径操作
	"strconv"                          //转换包
	"strings"
	"time" //时间类似Date
)

const (
	_DB_NAME        = "data/beeblog.db"
	_SQLITE3_DRIVER = "sqlite3"
)

type Category struct {
	Id              int64
	Title           string
	Ctype           int64
	Uid             int64  `orm:"null"`
	Created         string `orm:"index;null"`
	Views           int64  `orm:"index;null"`
	TopicTime       string `orm:"index;null"`
	TopicCount      int64  `orm:"null"`
	TopicLastUserId int64  `orm:"null"`
}

type Topic struct {
	Id              int64
	Uid             int64 `orm:"null"`
	Ctype           int64
	Category        string
	Labels          string
	Title           string
	Stick           int64
	Content         string `orm:"type(text)"`
	Attachment      string `orm:"null"`
	Created         string `orm:"index;null"`
	Updated         string `orm:"index;null"`
	Views           int64  `orm:"index;null"`
	Author          string `orm:"null"`
	ReplyTime       string `orm:"null"`
	ReplyCount      int64  `orm:"null"`
	ReplyLastUserId int64  `orm:"index;null"`
}

type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string `orm:"size(1000)"`
	Created string `orm:"index"`
}

type Gallery struct {
	Id      int64
	Name    string
	Uname   string
	Created string `orm:"index"`
}

type Speack struct {
	Id      int64
	Uname   string
	Content string `orm:"size(500);null"`
	Created string `orm:"index"`
}

type About struct {
	Id       int64
	Uname    string
	Nickname string
	Headpic  string
	Content  string `orm:"type(text)"`
	Created  string `orm:"index"`
}

func RegisterDB() {
	// if !com.IsExist(_DB_NAME) {
	// 	os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
	// 	os.Create(_DB_NAME)
	// }
	// orm.RegisterModel(new(Category), new(Topic), new(Comment))
	// orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	// orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
	orm.RegisterModel(new(Category), new(Topic), new(Comment), new(Gallery), new(Speack), new(About))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/myblog?charset=utf8")
}

func AddTopic(ctype, title, category, label, content, attachment string) error {
	//处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
	mtype, err := strconv.ParseInt(ctype, 10, 64)
	if err != nil {
		return err
	}
	//空格作为多个标签的分隔符

	o := orm.NewOrm()
	topic := &Topic{
		Title:      title,
		Ctype:      mtype,
		Category:   category,
		Labels:     label,
		Stick:      0,
		Content:    content,
		Attachment: attachment,
		Created:    time.Now().Local().Format("2006-01-02 15:04:05"),
		Updated:    time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	_, err = o.Insert(topic)
	if err != nil {
		return err
	}
	//更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		//如果存在就更新
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return err
}

func AddCategory(name, ctype string) error {
	c_type, err := strconv.ParseInt(ctype, 10, 4)

	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{
		Title: name,
		Ctype: c_type,
	}
	qs := o.QueryTable("category")
	err = qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}

func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64) //id 十进制 int64
	if err != nil {
		return err
	}

	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	o := orm.NewOrm()
	_, err = o.Insert(reply)

	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now().Local().Format("2006-01-02 15:04:05")
		topic.ReplyCount++
		_, err = o.Update(topic)
	}

	return err
}

//照片添加
func AddGallery(name, uname string) error {
	gallery := &Gallery{
		Name:    name,
		Uname:   uname,
		Created: time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	o := orm.NewOrm()
	_, err := o.Insert(gallery)

	if err != nil {
		return err
	}
	return nil
}

//照片添加
func AddSpeack(uname, content string) error {
	speack := &Speack{
		Uname:   uname,
		Content: content,
		Created: time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	o := orm.NewOrm()
	_, err := o.Insert(speack)

	if err != nil {
		return err
	}
	return nil
}

//照片添加
func AddAbout(uname, nickname, headpic, content string) error {
	about := &About{
		Uname:    uname,
		Nickname: nickname,
		Headpic:  headpic,
		Content:  content,
		Created:  time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	o := orm.NewOrm()
	_, err := o.Insert(about)

	if err != nil {
		return err
	}
	return nil
}

func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64) //id 十进制 int64
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return nil
}

func DeleteTopic(id string) error {
	tidNum, err := strconv.ParseInt(id, 10, 64) //id 十进制 int64
	if err != nil {
		return err
	}
	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}

	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}

	return err
}
func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64) //id 十进制 int64
	if err != nil {
		return err
	}
	var tidNum int64
	o := orm.NewOrm()
	reply := &Comment{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_, err = o.Delete(reply)
		if err != nil {
			return err
		}
	}
	replise := make([]*Comment, 0)
	/*
	   更新最后一个回复的时间
	*/
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replise)
	if err != nil {
		return err
	}
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		if len(replise) > 0 {
			topic.ReplyTime = replise[0].Created
		}

		topic.ReplyCount = int64(len(replise))
		_, err = o.Update(topic)
	}
	return err
}

func GetAllTopics(title, label, cate, sortFild, sortType string, pageSize, pageStart int64) (int64, []*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	var err error
	if len(cate) > 0 {
		qs = qs.Filter("category", cate)
	}
	if len(label) > 0 {
		qs = qs.Filter("labels__contains", "$"+label+"#")
	}
	var cnt int64
	if sortFild == "" || sortType == "" {
		_, err = qs.Filter("title__contains", title).Limit(pageSize, pageStart).All(&topics)
		cnt, err = qs.Filter("title__contains", title).Limit(pageSize, pageStart).Count()

	} else {
		_, err = qs.Filter("title__contains", title).OrderBy(sortType+sortFild).Limit(pageSize, pageStart).All(&topics)
		cnt, err = qs.Filter("title__contains", title).OrderBy(sortType+sortFild).Limit(pageSize, pageStart).Count()
	}
	for _, topic := range topics {
		topic.Labels = strings.Replace(strings.Replace(
			topic.Labels, "#", " ", -1), "$", "", -1)
	}

	return cnt, topics, err

}

func GetAllTopicsRand() ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	_, err := o.Raw("SELECT * FROM topic ORDER BY RAND() LIMIT 6").QueryRows(&topics)
	if err != nil {
		return nil, err
	}
	return topics, err

}
func GetAllTopicsByStick() ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	_, err := qs.Filter("stick", 1).Limit(8, 0).All(&topics)
	if err != nil {
		return nil, err
	}
	return topics, err

}
func DelGallery(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64) //id 十进制 int64
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	gallery := new(Gallery)

	qs := o.QueryTable("gallery")
	err = qs.Filter("id", cid).One(gallery)
	if err == nil {
		os.Remove(path.Join("attachment/"+gallery.Uname+"/gallery/300w", gallery.Name))
		os.Remove(path.Join("attachment/"+gallery.Uname+"/gallery", gallery.Name))
	}

	_, err = o.Delete(gallery)

	return nil
}

func GetAllCategory() ([]*Category, error) {
	o := orm.NewOrm()
	cate := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cate)
	return cate, err

}

func GetAllCategoryByType(ctype int64) ([]*Category, error) {
	o := orm.NewOrm()
	category := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.Filter("ctype", ctype).All(&category)
	if err != nil {
		return nil, err
	}
	return category, err
}

func GetAllReplies(tid string) ([]*Comment, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}
func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	topic.Views++
	_, err = o.Update(topic)

	topic.Labels = strings.Replace(strings.Replace(
		topic.Labels, "#", " ", -1), "$", "", -1)
	return topic, err
}

func ModifyTopic(tid, ctype, title, category, label, content, attachment string) error {
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	mtype, err := strconv.ParseInt(ctype, 10, 64)
	if err != nil {
		return err
	}

	var oldCate, oldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	beego.Error(topic)
	if o.Read(topic) == nil {
		oldCate = topic.Category
		beego.Error(oldCate)
		oldAttach = topic.Attachment
		topic.Ctype = mtype
		topic.Title = title
		topic.Labels = label
		if attachment != "" {
			topic.Attachment = attachment
		}
		topic.Category = category
		topic.Content = content
		topic.Updated = time.Now().Local().Format("2006-01-02 15:04:05")
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
		beego.Error(attachment)
	}
	//更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	if attachment != "" {
		//删除旧的附件
		if len(oldAttach) > 0 {
			os.Remove(path.Join("attachment", oldAttach))
		}
	}

	cate := new(Category)
	err = o.QueryTable("category").One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}

	return nil
}
func StickTopic(tid string) (int64, error) {

	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return 0, err
	}

	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldStick := topic.Stick
		if oldStick == 0 {
			topic.Stick = 1
		} else {
			topic.Stick = 0
		}
		_, err = o.Update(topic)
		if err != nil {
			return 0, err
		}
	}
	return topic.Stick, nil
}

func ModifyAbout(tid, uname, athor, attachment, content string) error {

	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	var oldAttach string
	o := orm.NewOrm()
	about := &About{Id: tidNum}
	beego.Info(about)
	if o.Read(about) == nil {
		oldAttach = about.Headpic
		beego.Info("原图片名--" + oldAttach)
		about.Nickname = athor
		about.Headpic = attachment
		about.Content = content
		_, err = o.Update(about)
		if err != nil {
			return err
		}
		beego.Info("现在的图片名" + attachment)
	}
	//删除旧的附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("/attachment/"+uname+"/about", oldAttach))
		beego.Info("删除原图片--" + oldAttach + "--成功！")
	}
	return nil
}
func GetAllGallery() ([]*Gallery, error) {
	o := orm.NewOrm()
	gallery := make([]*Gallery, 0)
	qs := o.QueryTable("gallery")
	_, err := qs.OrderBy("-created").All(&gallery)
	return gallery, err
}

func GetAllSpeack(pageSize, pageStart int64) (int64, []*Speack, error) {
	o := orm.NewOrm()
	speacks := make([]*Speack, 0)
	qs := o.QueryTable("speack")
	cnt, err := o.QueryTable("speack").Count()
	if err != nil {
		return cnt, nil, err
	}
	_, err = qs.OrderBy("-created").Limit(pageSize, pageStart).All(&speacks)
	return cnt, speacks, err
}
func GetAbout(uname string) (*About, error) {
	o := orm.NewOrm()
	about := new(About)
	qs := o.QueryTable("about")
	err := qs.Filter("uname", uname).One(about)
	if err != nil {
		return nil, err
	}
	return about, err
}
