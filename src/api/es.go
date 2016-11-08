package api

import (
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

type search struct {
}

func GetSearch() *search {
	return &search{}
}

func (s *search) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}
