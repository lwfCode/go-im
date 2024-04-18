package common

import "github.com/gin-gonic/gin"

func Response(c *gin.Context, code int, message string, reponse interface{}) {

	c.JSON(code, gin.H{
		"code": 200,
		"msg":  message,
		"data": reponse,
	})
}
