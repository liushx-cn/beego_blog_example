package controllers

import (
	"beego_blog_example/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

type CommentController struct {
	beego.Controller
}

type ClassController struct {
	BaseAuthController
}

type AuthController struct {
	BaseAuthController
}

// 获取文档列表，分类下的文档列表，以及搜索
func (c *ArticleController) List() {
	// 异步执行访客数据记录
	go models.AddNewRecord(c.Ctx.Request.RemoteAddr)

	logs.Info(c.Ctx.Request.RemoteAddr, " has visit article list on days ", strconv.Itoa(time.Now().Day()))

	p := c.GetString("page", "1")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}
	page--
	o := orm.NewOrm()
	var arts []*models.Article

	querySet := o.QueryTable(&models.Article{})
	total, _ := querySet.Count()
	_, e := querySet.RelatedSel().OrderBy("-create_time").Offset(page * 12).Limit(12).All(&arts)
	if e != nil {
		c.Data["json"] = map[string]string{
			"status":  "404",
			"message": "failed",
			"data":    e.Error(),
			"total":   string(total),
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"status":  "200",
			"message": "success",
			"data":    &arts,
			"total":   total,
		}
	}
	c.ServeJSON()
}

// 获取文档详情
func (c *ArticleController) Read() {
	pk := c.Ctx.Input.Param(":pk")
	art := new(models.Article)
	comments := []*models.Comment{}
	o := orm.NewOrm()
	querySet := o.QueryTable(art)
	e := querySet.Filter("id", pk).RelatedSel().One(art)
	o.QueryTable(new(models.Comment)).Filter("ArtObject", pk).All(&comments)
	if e == nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "200",
			"message": "success",
			"data": map[string]interface{}{
				"art":      art,
				"comments": comments,
			},
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"status":  "200",
			"message": "success",
			"data":    nil,
		}
	}
	c.ServeJSON()
}

// 创建文档
func (c *AuthController) Create() {
	art := new(models.Article)
	e := json.Unmarshal(c.Ctx.Input.RequestBody, art)

	cls := new(models.Classification)
	err := orm.NewOrm().QueryTable(cls).Filter("id", art.ClsId).One(cls)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "404",
			"message": "failed",
			"data":    "无效的分类",
		}
		c.ServeJSON()
		c.StopRun()
	}
	art.Class = cls

	data := ""
	if e != nil {
		data = e.Error()
	}
	art.Author = c.AuthUser
	o := orm.NewOrm()
	pk, err := o.Insert(art)
	if err != nil {
		data = err.Error()
	}
	data = strconv.Itoa(int(pk))
	c.Data["json"] = map[string]interface{}{
		"status":  "200",
		"message": "success",
		"data":    data,
	}
	c.ServeJSON()
}

// 更新文档
func (c *ArticleController) Update() {
	art := new(models.Article)
	e := json.Unmarshal(c.Ctx.Input.RequestBody, art)
	data := ""
	if e != nil {
		data = e.Error()
	}
	o := orm.NewOrm()
	num, err := o.Update(art, "title", "description", "Html", "code")
	if err != nil {
		data = err.Error()
	}
	data = string(num)
	c.Data["json"] = map[string]interface{}{
		"status":  "200",
		"message": "success",
		"data":    data,
	}
	c.ServeJSON()
}

// 检索文档
func (c *ArticleController) Search() {
	kw := c.GetString("search", "")
	arts := []*models.Article{}
	var total int64
	o := orm.NewOrm()
	querySet := o.QueryTable(new(models.Article))
	if kw != "" {
		total, _ = querySet.Filter("title_icontains", kw).OrderBy("-visit").All(&arts)
	} else {
		total, _ = querySet.OrderBy("-create_time").Offset(0).Limit(12).All(&arts)
	}
	c.Data["json"] = map[string]interface{}{
		"status":  "200",
		"message": "success",
		"data":    arts,
		"total":   total,
	}
	c.ServeJSON()
}

// 获取分类下的文档列表
func (c *ArticleController) Category() {
	cls := c.Ctx.Input.Param(":cls")
	p := c.GetString("page", "1")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}
	page--
	arts := []*models.Article{}

	o := orm.NewOrm()
	classification := new(models.Classification)
	querySet := o.QueryTable(classification)
	e := querySet.Filter("class", cls).RelatedSel().One(classification)
	if e == nil {
		arts = classification.Arts
	}
	c.Data["json"] = map[string]interface{}{
		"status":  "200",
		"message": "success",
		"data":    arts,
	}
	c.ServeJSON()
}

func (c *CommentController) Post() {
	comment := new(models.Comment)

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &comment)
	if err != nil {
		c.Data["json"] = map[string]string{
			"status":  "400",
			"message": "failed",
			"data":    err.Error(),
		}
		c.ServeJSON()
		c.StopRun()
	}
	artObj := new(models.Article)
	err2 := orm.NewOrm().QueryTable(new(models.Article)).Filter("id", comment.ArtId).One(artObj)
	if err2 != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "404",
			"message": "failed",
			"data":    "找不到该文档",
		}
		c.ServeJSON()
		c.StopRun()
	}
	comment.ArtObject = artObj

	o := orm.NewOrm()
	pk, e := o.Insert(comment)
	if e != nil {
		panic(e)
	}
	c.Data["json"] = map[string]interface{}{
		"status":  "200",
		"message": "success",
		"data":    pk,
	}
	c.ServeJSON()
}

func (c *ClassController) Post() {
	class := new(models.Classification)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, class)
	if err != nil {
		c.Data["json"] = map[string]string{
			"status":  "200",
			"message": "success",
			"data":    err.Error(),
		}
		c.ServeJSON()
	}
	o := orm.NewOrm()
	pk, e := o.Insert(class)
	if e != nil {
		panic(e)
	}
	c.Data["json"] = map[string]interface{}{
		"status":  "200",
		"message": "success",
		"data":    pk,
	}
	c.ServeJSON()
}
