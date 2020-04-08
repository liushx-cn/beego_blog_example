package models

import (
	"beego_blog_example/utils"
	"time"
)

type Account struct {
	BaseModel
	Username   string     `json:"username" orm:"size(64);unique;index"`
	Password   string     `json:"password" orm:"size(64)"`
	Nick       string     `json:"nick" orm:"size(64)"`
	Mobile     string     `json:"mobile" orm:"size(16);null;unique"`
	Avatar     string     `json:"avatar" orm:"size(128)"`
	Token      string     `json:"token" orm:"size(128)"`
	Permission int8       `json:"permission"` // 权限: 1, super 2, admin 3,normal 4, black
	Articles   []*Article `json:"article" orm:"reverse(many)"`
}

func (a *Account) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":         a.Id,
		"username":   a.Username,
		"nick":       a.Nick,
		"mobile":     a.Mobile,
		"avatar":     a.Avatar,
		"last_login": a.UpdateTime,
		"token":      a.Token,
	}
}

func (a *Account) CheckPassword(password string) bool {
	return utils.MD5(password) == a.Password
}

func (a *Account) GetPassword(password string) {
	a.Password = utils.MD5(password)
}

func (a *Account) GetToken() string {
	TokenStr := a.Username + string(time.Now().Unix())
	token := utils.MD5(TokenStr)
	a.Token = token
	return token
}
