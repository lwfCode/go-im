package router

import (
	"im/middlewares"
	"im/service"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	r := gin.Default()

	r.Use(cors.Default())

	/** 登陆接口 */
	r.POST("/api/login", service.Login)

	/** 注册接口 */
	r.POST("/api/register", service.Register)

	/** 发送验证码 */
	r.POST("/api/send", service.SendCode)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusNotFound,
			"msg":  "请求地址有误!",
		})
	})

	auth := r.Group("/api", middlewares.AuthCheck())
	{
		auth.GET("/user/details", service.GetUserDetails)
		auth.GET("/websocket/message", service.WebsocketMessage)
	}

	return r
}
