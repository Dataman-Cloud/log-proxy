package api

import (
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	STATUSBAD_TEST = "400-10002"
)

type monitor struct {
}

func GetMonitor() *monitor {
	return &monitor{}
}

func (m *monitor) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}
