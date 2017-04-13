package middleware

import (
	"strconv"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"github.com/gin-gonic/gin"
)

// Authenticate may add the authenticate methods.
func Authenticate(ctx *gin.Context) {
	ctx.Next()
}

// CORSMiddleware adds the CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		packPage(c)

		c.Next()
	}
}

func packPage(ctx *gin.Context) {
	page := models.Page{}
	if ctx.Query("from") != "" {
		if from, err := strconv.ParseInt(ctx.Query("from"), 10, 64); err == nil {
			page.RangeFrom = from
		} else {
			page.RangeFrom = ctx.Query("from")
		}
	} else {
		page.RangeFrom = nil
	}

	if ctx.Query("to") != "" {
		if to, err := strconv.ParseInt(ctx.Query("to"), 10, 64); err == nil {
			page.RangeTo = to
		} else {
			page.RangeTo = ctx.Query("to")
		}
	} else {
		page.RangeTo = nil
	}

	if size, err := strconv.Atoi(ctx.Query("size")); err == nil && size > 0 {
		page.PageSize = size
	} else {
		page.PageSize = 100
	}

	if p, err := strconv.Atoi(ctx.Query("page")); err == nil && p > 0 {
		page.PageFrom = (p - 1) * page.PageSize
	} else {
		page.PageFrom = 0
	}
	ctx.Set("page", page)
}
