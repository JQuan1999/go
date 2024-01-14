package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AuthSessionNotFound = 2001
)

func Auth(ctx *gin.Context) {
	sess := GetSession(ctx.GetHeader("Access-Token")) // 获取session
	if sess == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    AuthSessionNotFound,
			"message": "session not found or expired",
		})
		return
	}
	ctx.Next()
}
