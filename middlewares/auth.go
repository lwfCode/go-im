package middlewares

import (
	"im/common"
	"im/helper"

	"github.com/gin-gonic/gin"
)

const KEY = "token"

func AuthCheck() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		token := ctx.GetHeader(KEY)
		UserClaims, err := helper.AnalyseToken(token)
		if err != nil {
			ctx.Abort()
			common.Response(ctx, -1, "用户认证不通过", nil)
			return
		}
		ctx.Set("user_claims", UserClaims)
		ctx.Next()
	}
}
