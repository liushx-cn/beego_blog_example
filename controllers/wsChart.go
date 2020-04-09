package controllers

import (
	"container/list"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"time"
)

const (
	JOIN = iota
	TALK
	LEAVE
)

var (
	UserList = new(list.List)
	UserConn = make(map[string]*websocket.Conn)
	MsgPipe  = make(chan Msg, 10)
)

// 判断加入的用户是否是新用户
func Exists(username string) bool {
	for user := UserList.Front(); user != nil; user = user.Next() {
		if user.Value.(string) == username {
			return true
		}
	}
	return false
}

type Msg struct {
	User      string
	Type      int
	Content   string
	TimeStamp time.Time
}

func ChartCenter() {
	for {
		select {
		case msg := <-MsgPipe:
			message, _ := json.Marshal(msg)
			for _, v := range UserConn {
				if e := v.WriteMessage(websocket.TextMessage, message); e != nil {
					MsgPipe <- Msg{User: msg.User, Type: 3, Content: "leave", TimeStamp: time.Now()}
					// 从登录访客记录中移除当前访客， 并关闭连接
					delete(UserConn, msg.User)
				}
			}
		}
	}
}

type ChartController struct {
	beego.Controller
}

func (c *ChartController) Get() {
	nick := c.GetString("user")
	if nick == "" {
		nick = c.Ctx.Request.RemoteAddr
	}
	c.TplName = "websocket.html"
	c.Data["nick"] = nick
}
func (c *ChartController) Index() {
	c.TplName = "webChart.html"
}

func (c *ChartController) Join() {
	user := c.GetString("user")

	// 创建websocket对象
	weber := websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}
	// 将当前链接从HTTP转换为websocket
	conn, err := weber.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		panic(err)
	}
	if !Exists(user) {
		UserConn[user] = conn
		UserList.PushBack(user)
	}

	defer func() {
		MsgPipe <- Msg{User: user, Type: LEAVE, Content: "leave", TimeStamp: time.Now()}
		// 从登录访客记录中移除当前访客， 并关闭连接
		delete(UserConn, user)
		conn.Close()
		c.Data["json"] = map[string]string{"status": "200", "message": "over", "data": "会话结束"}
		c.ServeJSON()
	}()

	// 将新访客的消息加入消息队列
	msg := Msg{User: user, Type: JOIN, Content: "join", TimeStamp: time.Now()}
	MsgPipe <- msg

	// 循环监听
	for {
		_, p, e := conn.ReadMessage()
		if e != nil {
			logs.Error("== 通信终端，结束聊天" + e.Error())
			return
		}
		// 将获取的消息加入消息队列
		info := Msg{User: user, Type: TALK, Content: string(p), TimeStamp: time.Now()}
		MsgPipe <- info
	}
}

func init() {
	go ChartCenter()
}
