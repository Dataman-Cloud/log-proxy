package middleware

import (
	"github.com/gin-gonic/gin"
)

func Authenticate(ctx *gin.Context) {
	ctx.Next()
}
