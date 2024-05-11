package test

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8090", "http service address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

var ws = make(map[*websocket.Conn]struct{})

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader:", err)
		return
	}
	defer c.Close()
	ws[c] = struct{}{}
	for {
		/**======接收客户端发送的消息=========*/
		message := new(message)
		err := c.ReadJSON(message)
		if err != nil {
			log.Println("ReadJSON error:", err)
			return
		}
		fmt.Println(message.Data, "======消息message")
		fmt.Println(message.Type, "======消息type")

		/**======给客户端发消息=========*/
		for socket := range ws {
			err := socket.WriteMessage(websocket.TextMessage, []byte("hello world"))
			if err != nil {
				log.Println("消息发送失败~~")
				return
			}
		}
	}
}

func TestWebSocket(t *testing.T) {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
