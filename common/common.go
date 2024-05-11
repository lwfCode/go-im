package common

import (
	"crypto/md5"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
