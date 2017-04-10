package api

import (
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

// Monitor struct
type Monitor struct {
}

// NewMonitor init the struct monitor
func NewMonitor() *Monitor {
	return &Monitor{}
}

// Ping return string "success"
func (m *Monitor) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}
