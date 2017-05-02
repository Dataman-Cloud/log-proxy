package api

import (
	"fmt"
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

// Query return the results of quering from prometheus
func (m *Monitor) Query(ctx *gin.Context) {
	param := &models.QueryParameter{
		Metric: ctx.Query("metric"),
		App:    ctx.Query("app"),
		Task:   ctx.Query("task"),
		Start:  ctx.Query("start"),
		End:    ctx.Query("end"),
		Step:   ctx.Query("step"),
		Expr:   ctx.Query("expr"),
	}

	if param.Metric != "" && param.Expr != "" {
		err := fmt.Errorf("The paramter confict between metric and expr")
		utils.ErrorResponse(ctx, err)
		return
	}

	if param.Metric == "" && param.Expr == "" {
		err := fmt.Errorf("The paramter metric or expr required")
		utils.ErrorResponse(ctx, err)
		return
	}

	if param.Metric != "" && param.App == "" {
		err := fmt.Errorf("The paramter metric and app required")
		utils.ErrorResponse(ctx, err)
		return
	}
	if param.Expr != "" {
		query := &service.Query{
			ExprTmpl:       service.SetQueryExprsList(),
			HTTPClient:     http.DefaultClient,
			PromServer:     config.GetConfig().PrometheusURL,
			Path:           service.QUERYRANGEPATH,
			QueryParameter: param,
		}
		data, err := query.QueryExpr()
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
		utils.Ok(ctx, data)
		return
	}

	if param.Metric != "" {
		query := &service.Query{
			ExprTmpl:       service.SetQueryExprsList(),
			HTTPClient:     http.DefaultClient,
			PromServer:     config.GetConfig().PrometheusURL,
			Path:           service.QUERYRANGEPATH,
			QueryParameter: param,
		}
		data, err := query.QueryMetric()
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
		utils.Ok(ctx, data)
		return
	}
}

// GetApps return the items of query metrics
func (m *Monitor) GetApps(ctx *gin.Context) {
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

	data, err := query.GetQueryApps()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// GetAppsTasks return the items of query metrics
func (m *Monitor) GetAppsTasks(ctx *gin.Context) {
	param := &models.QueryParameter{
		Start:   ctx.Query("start"),
		End:     ctx.Query("end"),
		Step:    ctx.Query("step"),
		Cluster: ctx.Param("clusterid"),
		App:     ctx.Param("appid"),
	}

	query := &service.Query{
		ExprTmpl:       service.SetQueryExprsList(),
		HTTPClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PrometheusURL,
		Path:           service.QUERYRANGEPATH,
		QueryParameter: param,
	}

	data, err := query.GetQueryAppTasks()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}
