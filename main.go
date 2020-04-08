package main

import (
	"beego_blog_example/models"
	_ "beego_blog_example/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 创建超级用户使用的函数
func CreateSuper() {
	user := models.Account{
		Username: "haiton",
		Nick:     "Haiton",
		Mobile:   "18355556666",
	}
	user.GetPassword("superPassword")
	o := orm.NewOrm()
	pk, _ := o.Insert(&user)
	println(pk)
}

func init() {
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", "root:mysql@tcp(127.0.0.1:3306)/myblog?charset=utf8mb4")

	// https://beego.me/docs/module/logs.md
	// 配置日志，单文件，设定存入级别
	//_ = logs.SetLogger(logs.AdapterFile, `{"filename":"D:\\AppData\\go\\src\\Haiton\\logs\\error.log","level":3,"daily":true,"maxdays":10,"rotate":true}`)
	// 多文件日志，生产环境一般都会采用这种方式， filename设置的文件通用部分，如下，级别设置alert error等，则相应的级别的日志信息会分别创建 web.alert.log web.error.log 存储
	_ = logs.SetLogger(
		logs.AdapterMultiFile,
		`{"filename":"D:\\AppData\\go\\src\\Haiton\\logs\\web.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],"maxdays":10}`)

	// 配置日志监控提醒，同样可以配置级别，级别及以上的信息会发送邮件提醒
	//_ = logs.SetLogger(logs.AdapterMail, `{"username":"example@gmail.com","password":"password","host":"smtp.gmail.com:587","sendTos":["someone@gmail.com"]}`)
	//_ = beego.BeeLogger.DelLogger(logs.AdapterConsole)
	logs.EnableFuncCallDepth(true)
}

func main() {
	o := orm.NewOrm()
	_ = o.Using("default")
	orm.RunCommand()
	beego.Run()
}
