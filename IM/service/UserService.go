package service

import (
	"IM/models"
	"IM/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)
import "IM/sql"

type UserService struct {
	sql.Userdao
}

func (T *UserService) Router(g *gin.Engine) {
	group := g.Group("/user")
	group.POST("register", T.register)
	group.POST("login", T.Login)
	//group.GET("sendmsg", SendMsg)
	group.GET("senduser", SendUsermsg)
	group.GET("usermsg", T.usermsg)
	group.GET("test", T.Test)
}
func (T *UserService) Test(con *gin.Context) {
	msg := models.Message{TargetId: 1}
	con.JSON(200, &msg)
}
func (T *UserService) usermsg(con *gin.Context) {
	userId := con.Query("userId")
	targetId := con.Query("targetId")
	revRange := make([]string, 0)
	revRange = nil
	if userId < targetId {
		revRange = sql.RedZRevRange("msg:"+userId+"_"+targetId, 0, -1)
	} else {
		revRange = sql.RedZRevRange("msg:"+targetId+"_"+userId, 0, -1)
	}
	if revRange != nil {
		con.JSON(200, &revRange)
		return
	}
	con.JSON(200, "出错了")
}

func (T *UserService) register(con *gin.Context) {
	var newu struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	con.BindJSON(&newu)

	if T.Find(models.UserBase{
		Name:     newu.Name,
		Password: newu.Password,
	}).Name == newu.Name {
		con.JSON(200, "失败，账号已存在！")
		return
	}

	salt := fmt.Sprintf("%d", rand.Int()%1000)
	code := utils.MD5Encode(newu.Password, salt)
	base := models.UserBase{Name: newu.Name, Password: code, Salt: salt}
	T.Insert(base)
	con.JSON(200, "成功插入")
	return
}

func (T *UserService) Login(con *gin.Context) {
	var newu struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	con.BindJSON(&newu)
	if sql.Redget(newu.Name) == newu.Password {
		con.JSON(200, "redis登录成功")
		return
	}
	u := T.Find(models.UserBase{
		Name: newu.Name,
		//Password: newu.Password,
	})

	if utils.Md5Decode(newu.Password, u.Salt, u.Password) {
		redstr := sql.Redstr(newu.Name, newu.Password)
		log.Println(redstr)
		con.JSON(200, "登录成功")
		return
	}
	con.JSON(200, "密码错误或账号不存在")
	return

}

var up = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err2 := ws.Close()
		if err != nil {
			log.Println(err2)
		}
	}(ws)

	MsgHandler(ws, c)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		subscribe, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			log.Println(err)
		}
		tm := time.Now().Format("2002-01-02 15:11:11")
		m := fmt.Sprintf("[ws][%s]", tm, subscribe)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUsermsg(c *gin.Context) { //可能冒充身份
	Chat(c.Writer, c.Request)
}
