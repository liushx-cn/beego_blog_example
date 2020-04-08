package routers

import (
	"beego_blog_example/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 获取首页文档列表
	beego.Router("/arts", &controllers.ArticleController{}, "get:List")

	// 获取分类下的文档列表
	beego.Router("/arts/:cls(\\d{1,6})", &controllers.ArticleController{}, "get:Category")

	// 创建文档
	beego.Router("/art", &controllers.AuthController{}, "post:Create")

	// 获取文档详情或更新文档
	beego.Router("/art/:pk(\\d{1,6})", &controllers.ArticleController{}, "get:Read;put:Update")

	// 搜索,从标题中搜索相关文档
	beego.Router("/search", &controllers.ArticleController{}, "get:Search")

	// 提交评论
	beego.Router("/comment", &controllers.CommentController{})

	// 创建分类
	beego.Router("/class", &controllers.ClassController{})
}
