package service

import (
	"IM/models"
	"IM/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

//消息

type Node struct {
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息
	GroupSets     set.Interface   //好友 / 群
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	userIdstring := query.Get("userId")
	userId, _ := strconv.ParseInt(userIdstring, 10, 64)
	//token := query.Get("token")
	//msgtype := query.Get("type")
	//	targetId := query.Get("targetId")
	//	context := query.Get("context")
	isvalida := true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.NonThreadSafe),
	}
	//用户关系
	//
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	u := sql.Userdao{}.FindById(userId)
	sql.Redstr("userId:"+string(userId), u)
	go sendProc(node)
	go recvProc(node)
	//	sendMsg(userId, []byte("wolcome"+strconv.Itoa(int(userId))))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("发送了一条信息：" + string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data)
		//broadMsg(data)
		//	node.DataQueue <- data
		//	fmt.Println("[ws]<<<<<<", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}
func ini1t() {
	go udpSendProc()
	go updRecvProc()
}

func dispatch(data []byte) {
	msg := models.Message{}
	err := json.Unmarshal(data, &msg)
	data, _ = json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
	case 2: //群发
		sendGroup(msg.TargetId, msg.UserId, data)
		//case 3://广播
		//	sendAllmsg()
	}
}
func sendMsg(targetId int64, msg []byte) {
	jsonMsg := models.Message{}
	json.Unmarshal(msg, &jsonMsg)
	key := ""
	x := strconv.FormatInt(targetId, 10)
	y := strconv.FormatInt(jsonMsg.UserId, 10)

	if x < y {
		key = "msg:" + x + "_" + y
	} else {
		key = "msg:" + y + "_" + x
	}
	res := sql.RedZRevRange(key, 0, -1)
	score := float64(cap(res)) + 1

	sql.RedZAdd(key, score, msg)
	rwLocker.RLock()
	node, ok := clientMap[targetId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
func sendGroup(targetId int64, userId int64, msg []byte) {
	t := make([]int64, 0)
	redget := sql.RedListGetAll("Group:" + strconv.FormatInt(targetId, 10))
	if len(redget) != 0 {
		for i := 0; i < len(redget); i++ {
			if atoi, _ := strconv.Atoi(redget[i]); int64(atoi) != userId {
				//不能给自己发消息
				sendMsg(int64(atoi), msg)
			}
		}
		return
	}
	contactdao := sql.Contactdao{}
	find3 := contactdao.Find2(targetId)
	for i := 0; i < len(find3); i++ {
		t = append(t, int64(find3[i].TargetId))
		if int64(find3[i].TargetId) != userId { //不能给自己发消息
			sendMsg(int64(find3[i].TargetId), msg)
		}
	}
	sql.RedListGroup("Group:"+strconv.FormatInt(targetId, 10), t)
}

func udpSendProc() {
	udp, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(172, 22, 128, 1),
		Port: 3000,
	})
	defer udp.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpsendChan:
			_, err := udp.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func updRecvProc() {
	udp, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		log.Println(err)
	}
	defer udp.Close()
	for {
		var buf [512]byte
		n, err := udp.Read(buf[0:])
		if err != nil {
			log.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}
