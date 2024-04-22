package service

import (
	"im/common"
	"im/helper"
	"im/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var ws = make(map[string]*websocket.Conn, 0)

type Message struct {
	RoomIdentity string `json:"room_identity"`
	Data         string `json:"data"`
}

func WebsocketMessage(c *gin.Context) {
	cli, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.Response(c, http.StatusBadRequest, "连接websocket异常", err)
		return
	}
	defer cli.Close()

	user := c.MustGet("user_claims").(*helper.UserClaims)

	//多个用户连接放到map中
	ws[user.Identity] = cli

	for {
		msg := new(Message)
		err = cli.ReadJSON(msg)
		if err != nil {
			log.Printf("read Error:%v\n", err)
			return
		}
		_, err := models.GetUserRoomByUserIdentityRoomIdentity(user.Identity, msg.RoomIdentity)
		if err != nil {
			log.Printf("userIdentity:%s and roomIdentity:%s : Not Exits\n", user.Identity, msg.RoomIdentity)
			return
		}
		//保存消息
		formattedTime := time.Now().Format("2006-01-02 15:04:05")
		mb := &models.MessageBasic{
			UserIdentity: user.Identity,
			RoomIdentity: msg.RoomIdentity,
			Data:         msg.Data,
			CreatedAt:    formattedTime,
			UpdatedAt:    formattedTime,
		}
		err = models.InsertOneMessageBasic(mb)
		if err != nil {
			log.Printf("保存消息失败,[ERROR]:%v\n", err)
			return
		}
		userRoomIds, err := models.GetUserRoomByRoomIdentity(msg.RoomIdentity)
		if err != nil {
			log.Printf("GetUserRoomByRoomIdentity Error:%v\n", err)
			return
		}
		for _, v := range userRoomIds {
			if socket, ok := ws[v.UserIdentity]; ok {
				err := socket.WriteMessage(websocket.TextMessage, []byte(msg.Data))
				if err != nil {
					log.Printf("WriteMessage Error:%v\n", err)
					return
				}
			}
		}
	}
}
