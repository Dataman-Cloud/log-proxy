package api

import (
	"github.com/Dataman-Cloud/log-proxy/src/service"
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

// GetQueryItems return the items of query metrics
func (m *Monitor) GetQueryItems(ctx *gin.Context) {
	query := &service.Query{
		ExprTmpl: service.SetQueryExprsList(),
	}
	utils.Ok(ctx, query.GetQueryItemList())
}
