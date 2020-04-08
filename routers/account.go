package routers

import (
	"beego_blog_example/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/user", &controllers.UserController{})
}
