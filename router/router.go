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
	r.POST("/login", service.Login)

	/** 注册接口 */
	r.POST("/register", service.Register)

	/** 发送验证码 */
	r.POST("/send", service.SendCode)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusNotFound,
			"msg":  "请求地址有误!",
		})
	})

	userGroup := r.Group("/v1/api", middlewares.AuthCheck())
	{
		userGroup.GET("/user/details", service.GetUserDetails)
	}

	return r
}
