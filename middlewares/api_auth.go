package middlewares

import (
	"ginchat/common"
	"ginchat/define"
	"ginchat/helper"

	"github.com/gin-gonic/gin"
)

const USER_CLAIMS = "user_claims"

func ApiAuth() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		token := ctx.GetHeader(define.V1API_HEADER_TOEKN)
		UserClaims, err := helper.AnalyseToken(token)
		if err != nil {
			ctx.Abort()
			common.Response(ctx, 403, "用户认证不通过", nil)
			return
		}
		ctx.Set(USER_CLAIMS, UserClaims)
		ctx.Next()
	}
}
