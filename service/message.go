package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

var clientManager = NewClientManager()

func SendChat(c *gin.Context) {

	// user, err := CheckAuth(c)
	// if err != nil {
	// 	log.Println("获取失败，token校验失败!")
	// 	return
	// }

	userId, _ := strconv.Atoi(c.Query("user_id")) //登录的用户
	// tarUid, _ := strconv.Atoi(c.Query("tar_id"))  //向谁发起聊天
	// fmt.Println(tarUid, "======接收者id")
	clientManager.start(userId, c)

}
