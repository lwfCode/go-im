package router

import (
	"ginchat/docs"
	"ginchat/middlewares"
	"ginchat/service"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {

	r := gin.Default()

	docs.SwaggerInfo.BasePath = ""

	r.Use(cors.Default())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusNotFound,
			"msg":  "请求地址有误!",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	/**========================v1版本接口 start====================== */
	v1Verison := "/v1/api"
	/** 无需授权 */
	noAuth := r.Group(v1Verison)
	{
		/** 发送验证码 */
		//noAuth.POST("/send", service.SendCode)

		/** 注册接口 */
		noAuth.POST("/user/register", service.Register)

		/** 登陆接口 */
		noAuth.POST("/user/login", service.Login)

		//文件上传
		noAuth.POST("/file/upload", service.FileUpload)

		noAuth.GET("/user/sendMsg", service.SendChat)
	}

	apiAuth := r.Group(v1Verison, middlewares.ApiAuth())
	{
		apiAuth.GET("/user/userList", service.GetUserList)
		apiAuth.GET("/user/details", service.GetUserDetails)
		apiAuth.GET("/user/info", service.GetUserInfo)
		apiAuth.POST("/user/delete", service.DeleteUser)
		apiAuth.POST("/user/update", service.UpdateUser)

		//获取缓存在redis中的消息
		apiAuth.GET("/user/messageList", service.MessageList)

		//发送消息
		// apiAuth.GET("/user/sendMsg", service.SendChat)
	}
	/** ========================v1版本接口 End====================== */

	// noAuth.GET("/api/getUserList", service.GetUserList)

	return r
}
