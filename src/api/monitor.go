package api

import (
	"net/http"

	"github.com/Dataman-Cloud/log-proxy/src/service"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	STATUSBADTEST = "400-10002"
)

const (
	QUERYRANGEPATH = "/api/v1/query_range"
	QUERYPATH      = "/api/v1/query"
)

type monitor struct {
}

func NewMonitor() *monitor {
	return &monitor{}
}

func (m *monitor) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

func (m *monitor) QueryRange(ctx *gin.Context) {
	query := &service.QueryRange{
		HttpClient: http.DefaultClient,
		PromServer: config.GetConfig().PROMETHEUS_URL,
		Path:       QUERYRANGEPATH,
		Type:       ctx.Query("type"),
		AppID:      ctx.Query("appid"),
		TaskID:     ctx.Query("taskid"),
		Metric:     ctx.Query("metric"),
		From:       ctx.Query("from"),
		To:         ctx.Query("to"),
		Step:       ctx.Query("step"),
	}
	data := service.NewMetricList()
	err := data.GetMetricList(query)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

func (m *monitor) QueryApps(ctx *gin.Context) {
	query := &service.QueryRange{
		HttpClient: http.DefaultClient,
		PromServer: config.GetConfig().PROMETHEUS_URL,
		Path:       QUERYPATH,
		AppID:      ctx.Query("appid"),
	}

	apps := service.NewAppsList()
	err := apps.GetAppsList(query)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, apps)
}

func (m *monitor) QueryNodes(ctx *gin.Context) {
	query := &service.QueryRange{
		HttpClient: http.DefaultClient,
		PromServer: config.GetConfig().PROMETHEUS_URL,
		Path:       QUERYRANGEPATH,
		Node:       ctx.Query("node"),
	}

	data := service.NewNodesMetric()
	err := data.GetNodesMetric(query)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

func (m *monitor) QueryApp(ctx *gin.Context) {
	query := &service.QueryRange{
		HttpClient: http.DefaultClient,
		PromServer: config.GetConfig().PROMETHEUS_URL,
		Path:       QUERYRANGEPATH,
		AppID:      ctx.Query("appid"),
	}

	data := service.NewAppMetric()
	err := data.GetAppMetric(query)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}
