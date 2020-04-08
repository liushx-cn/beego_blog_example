package controllers

import (
	"beego_blog_example/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BaseAuthController struct {
	beego.Controller
	AuthUser *models.Account
}

func (c *BaseAuthController) Prepare() {
	authStr := c.Ctx.Input.Header("Authenticate")
	if authStr == "" {
		c.Data["json"] = map[string]string{"status": "401", "message": "invalid Authentication", "data": "Authenticate 为空，无法登陆"}
		c.ServeJSON()
		c.StopRun()
	}

	account := new(models.Account)
	o := orm.NewOrm()

	err := o.QueryTable(account).Filter("token", authStr).One(account)
	switch err {
	case orm.ErrMultiRows:
		c.Data["json"] = map[string]string{"status": "401", "message": "no auth", "data": "没有登录"}
		c.ServeJSON()
		c.StopRun()
	case orm.ErrNoRows:
		c.Data["json"] = map[string]string{"status": "401", "message": "no auth", "data": "没有登录"}
		c.ServeJSON()
		c.StopRun()
	}
	c.AuthUser = account
}
