package api

import (
	"net/http"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
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

// GetUsers return the items of query metrics
func (m *Monitor) GetUsers(ctx *gin.Context) {
	param := &models.QueryParameter{
		Start: ctx.Query("start"),
		End:   ctx.Query("end"),
		Step:  ctx.Query("step"),
	}

	query := &service.Query{
		ExprTmpl:       service.SetQueryExprsList(),
		HTTPClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PrometheusURL,
		Path:           service.QUERYRANGEPATH,
		QueryParameter: param,
	}

	data, err := query.GetQueryUsers()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}
