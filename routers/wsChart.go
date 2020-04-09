package routers

import (
	"beego_blog_example/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 先进入聊天首页，输入昵称
	beego.Router("/chart", &controllers.ChartController{}, "get:Index")
	// 进入聊天页，自动建立连接
	beego.Router("/wschart", &controllers.ChartController{})
	// 建立连接，并开始聊天
	beego.Router("/ws/join", &controllers.ChartController{}, "get:Join")
}
