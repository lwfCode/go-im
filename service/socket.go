package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/define"
	"ginchat/utils"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

type Client struct {
	Socket        *websocket.Conn
	Addr          string        //客户端地址
	FirstTime     uint64        //首次连接时间
	UserId        int           //用户id
	HeartbeatTime uint64        //用户上次心跳时间
	LoginTime     uint64        //登陆时间
	SendData      chan []byte   //消息
	GroupSets     set.Interface // 好友/群
}

// 初始化连接管理器
func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[int]*Client, 1000),
		Register:   make(chan *Client, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}
	return
}

// 初始化websocket连接
func NewClient(addr string, cli *websocket.Conn, userId int, firstTime uint64) (client *Client) {

	client = &Client{
		Addr:          addr,
		Socket:        cli,
		UserId:        userId,
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
		SendData:      make(chan []byte, 100),
	}
	return client
}

// 启动websocket
func (manager *ClientManager) start(userId int, c *gin.Context) {
	cli, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {

		return true
	}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err, "+++++ERR")
		http.NotFound(c.Writer, c.Request)
		return
	}
	currentTime := uint64(time.Now().Unix())
	client := NewClient(cli.RemoteAddr().String(), cli, userId, currentTime)

	go manager.HanldManager()

	go client.read()
	go client.write()
	manager.Register <- client

	fmt.Println("websocet 程序启动成功")
}

// 读取客户端消息
func (c *Client) read() {
	defer func() {
		fmt.Println("读取客户端数据 关闭send..", c)
		close(c.SendData)
	}()

	defer func() {
		clientManager.Unregister <- c

		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Println("ReadMessage recv[ERROR]:", err)
			return
		}
		fmt.Printf("read 读取客户端消息内容===>%v\n", string(message))

		//处理客户端发送的数据
		dispatch(c, message)
	}
}

// 往客户端发送消息
func (client *Client) write() {

	defer func() {
		client.Socket.Close()
		fmt.Println("Client发送数据 关闭", client)
	}()

	for {
		select {
		case message := <-client.SendData:
			// if !ok {
			// 	//发送数据失败，关闭连接
			// 	fmt.Println("发送数据失败，关闭连接====")
			// 	return
			// }
			err := client.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("往客户端发送数据失败=====>", err)
				return
			}
		default:
			time.Sleep(time.Millisecond * 2)
		}

	}
}

// 处理数据
func dispatch(client *Client, message []byte) {

	request := &define.SendMsg{}
	// request.Date = time.Now().Format("2006-01-02 15:04:05")
	request.Date = time.Now().Format("2006年01月02日 15:04:05")

	err := json.Unmarshal(message, request)
	if err != nil {
		fmt.Println("处理数据 json Unmarshal", err)
		// client.SendMsg([]byte("数据不合法"))
		return
	}
	requestData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		client.sendErr([]byte("处理数据失败"))
		return
	}
	cmd := request.Cmd

	switch cmd {
	case "ping":
		fmt.Println("心跳=====>", cmd)
	case "chat": //单聊
		// fmt.Println("requestData ===> ", string(requestData))
		client.SendMsg(requestData, request.TargetId)
	default:
		client.sendErr([]byte("数据不合法"))
		return
	}
}

func (client *Client) SendMsg(message []byte, TargetId int) {
	if client == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()
	err, cli := clientManager.GetUserKeyClient(TargetId)
	if err != nil {
		fmt.Println(err, "+++++++++GetUserKeyClient")
		return
	} else if cli != nil {
		cli.SendData <- message
	} else {
		// client.SendData <- message
	}
	var key string
	targetIdStr := strconv.Itoa(TargetId)
	userIdStr := strconv.Itoa(client.UserId)

	if TargetId > client.UserId {
		key = "chat_" + userIdStr + "_" + targetIdStr
	} else {
		key = "chat_" + targetIdStr + "_" + userIdStr
	}
	ctx := context.Background()
	resultSplice, err := utils.Redis.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err, "+++++redis.ZRevRange ERROR")
		return
	}
	fmt.Println(resultSplice, "++++resultSplice")
	score := float64(cap(resultSplice)) + 1

	_, err = utils.Redis.ZAdd(ctx, key, &redis.Z{score, message}).Result()
	if err != nil {
		fmt.Println(err, "++++ZADD ERROR")
		return
	}
}

func (client *Client) sendErr(message []byte) {
	if client == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()
	client.SendData <- message
}
