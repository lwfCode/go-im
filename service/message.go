package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

var clientManager = NewClientManager()

func SendChat(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id")) //登录的用户
	clientManager.start(userId, c)
}
