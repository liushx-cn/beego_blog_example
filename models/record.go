package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// 访客记录
type Visitor struct {
	BaseModel
	IpAddr string `json:"ip_addr" orm:"size(16)"`
}

func AddNewRecord(IP string) {
	visit := Visitor{IpAddr: IP}
	o := orm.NewOrm()
	_, e := o.Insert(&visit)
	if e != nil {
		logs.Error(e.Error())
	}
}
