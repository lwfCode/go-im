package service

import (
	"fmt"
	"sync"
	"time"
)

type ClientManager struct {
	Clients     map[*Client]bool // 全部的连接
	ClientsLock sync.RWMutex     // 读写锁
	Users       map[int]*Client  // 登录的用户 // userId
	UserLock    sync.RWMutex     // 读写锁
	Register    chan *Client     // 连接连接处理
	Unregister  chan *Client     // 断开连接处理程序
	Broadcast   chan []byte      // 广播 向全部成员发送数据
}

// 处理管道
func (manager *ClientManager) HanldManager() {
	for {
		select {
		case client := <-manager.Register:
			//建立连接
			manager.EventRegister(client)
		case client := <-manager.Unregister:
			//断开连接
			manager.EventClose(client)
		case data := <-manager.Broadcast:
			manager.EventBroad(data)
		default:
			time.Sleep(time.Millisecond * 2)
		}

	}
}

// 广播消息
func (manager *ClientManager) EventBroad(data []byte) {
	keys := manager.GetUserKeys()

	fmt.Println("所有用户的key======>", keys)
	fmt.Println("广播的数据=======>", string(data))
}

// 建立连接客户端
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)
	manager.AddUsers(client)

	broadData := []byte("您好，欢迎进入IM聊天系统!")
	manager.Broadcast <- broadData

	fmt.Println("EventRegister 用户建立连接", client.Addr)
}

// 断开链接
func (manager *ClientManager) EventClose(client *Client) {
	manager.DelClients(client)

	// 删除用户连接
	deleteResult := manager.DelUsers(client)
	if deleteResult == false {
		// 不是当前连接的客户端
		return
	}

	fmt.Println("EventUnregister 用户断开连接", client.Addr)
}

// 添加用户
func (manager *ClientManager) AddUsers(client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	manager.Users[client.UserId] = client
}

func (manager *ClientManager) DelUsers(client *Client) (result bool) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()
	if value, ok := manager.Users[client.UserId]; ok {
		if value.Addr != client.Addr {
			return
		}
		delete(manager.Users, client.UserId)
		result = true
	}

	return
}

// 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	manager.Clients[client] = true
	manager.ClientsLock.Unlock()
}

// 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		fmt.Println("断开链接～～～～")
		delete(manager.Clients, client)
	}
}

// 获取所有连接的用户ids
func (manager *ClientManager) GetUserKeys() (userKeys []int) {

	userKeys = make([]int, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	for key, _ := range manager.Users {
		userKeys = append(userKeys, key)
	}

	return
}

/** 获取单个用户的连接*/
func (manager *ClientManager) GetUserKeyClient(userId int) (err error, client *Client) {
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	client, ok := manager.Users[userId]
	if !ok {
		return err, nil
	}
	return nil, client
}
