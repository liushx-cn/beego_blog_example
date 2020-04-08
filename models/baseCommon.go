package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type BaseModel struct {
	Id         int64     `json:"id"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `json:"update_time" orm:"auto_now;type(datetime)"`
}

//func (b *BaseModel) Create() int64 {
//	o := orm.NewOrm()
//	pk, err := o.Insert(b)
//	if err != nil {
//		panic(err)
//	}
//	return pk
//}
//
//func (b *BaseModel) BulkInsert(blk int, objs []*BaseModel) int {
//	if blk != len(objs) {
//		panic(errors.New("批量插入的实际对象数量和声明的不一致"))
//	}
//	o := orm.NewOrm()
//	successNum, err := o.InsertMulti(blk, objs)
//	if err != nil {
//		panic(err)
//	}
//	return int(successNum)
//}
//
//func (b *BaseModel) GetQuerySet() orm.QuerySeter {
//	o := orm.NewOrm()
//	querySet := o.QueryTable(b)
//	return querySet
//}

func init() {
	orm.RegisterModel(new(Account), new(Article), new(Classification), new(Comment), new(Visitor))
}
