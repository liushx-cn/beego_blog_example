package controllers

import (
	"beego_blog_example/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}
type UserController struct {
	BaseAuthController
}

func (c *LoginController) Post() {
	data := map[string]string{"username": "", "password": ""}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &data)
	if err != nil {
		panic(err)
	}
	o := orm.NewOrm()
	user := new(models.Account)
	user.Username = data["username"]

	querySet := o.QueryTable(user)
	errDb := querySet.Filter("username", user.Username).One(user)
	if errDb != nil {
		c.Data["json"] = map[string]string{"status": "400", "message": "no account", "data": "该用户不存在"}
		c.ServeJSON()
		c.StopRun()
	}
	if user.CheckPassword(data["password"]) {
		user.GetToken()
		_, e := o.Update(user, "token")
		if e != nil {
			panic(e)
		}
		c.Data["json"] = map[string]interface{}{"status": "200", "message": "success", "data": user.GetInfo()}
	} else {
		c.Data["json"] = map[string]string{"status": "400", "message": "wrong password", "data": "密码错误"}
	}
	c.ServeJSON()
}

func (c UserController) Get() {
	user := c.AuthUser
	c.Data["json"] = map[string]interface{}{"status": "200", "message": "success", "data": user.GetInfo()}
	c.ServeJSON()
}

func (c *UserController) Put() {
	user := c.AuthUser
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		panic(err)
	}
	o := orm.NewOrm()
	_, e := o.Update(user, "nick", "avatar")
	if e != nil {
		panic(e)
	}
	c.Data["json"] = map[string]interface{}{"status": "200", "message": "success", "data": user.GetInfo()}
	c.ServeJSON()
}
