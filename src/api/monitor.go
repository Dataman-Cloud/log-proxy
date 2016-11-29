package api

import (
	"net/http"

	"github.com/Dataman-Cloud/log-proxy/src/backends"
	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	STATUSBADTEST = "400-10002"
)

const (
	QUERYRANGEPATH = "/api/v1/query_range"
	QUERYPATH      = "/api/v1/query"
	RULESPATH      = "/api/v1/rules"
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

// Backends API: PromQL
func (m *monitor) PromqlQuery(ctx *gin.Context) {
	promql := &backends.Promql{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().PROMETHEUS_URL,
		Path:       QUERYPATH,
		Query:      ctx.Query("query"),
		Time:       ctx.Query("time"),
	}

	data, err := promql.GetPromqlQuery()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (m *monitor) PromqlQueryRange(ctx *gin.Context) {
	promql := &backends.Promql{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().PROMETHEUS_URL,
		Path:       QUERYRANGEPATH,
		Query:      ctx.Query("query"),
		Start:      ctx.Query("start"),
		End:        ctx.Query("end"),
		Step:       ctx.Query("step"),
	}

	data, err := promql.GetPromqlQuery()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// Backends API: AlertManager
const (
	ALERTSPATH       = "/api/v1/alerts"
	ALERTSGROUSPPATH = "/api/v1/alerts/groups"
	ALERTSSTATUSPATH = "/api/v1/status"
)

func (m *monitor) GetAlerts(ctx *gin.Context) {
	query := &backends.AlertManager{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().ALERTMANAGER_URL,
		Path:       ALERTSPATH,
	}

	data, err := query.GetAlertManagerResponse()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (m *monitor) GetAlertsGroups(ctx *gin.Context) {
	query := &backends.AlertManager{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().ALERTMANAGER_URL,
		Path:       ALERTSGROUSPPATH,
	}

	data, err := query.GetAlertManagerResponse()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (m *monitor) GetAlertsStatus(ctx *gin.Context) {
	query := &backends.AlertManager{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().ALERTMANAGER_URL,
		Path:       ALERTSSTATUSPATH,
	}

	data, err := query.GetAlertManagerResponse()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (m *monitor) GetAlertsRules(ctx *gin.Context) {
	query := &backends.AlertManager{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().PROMETHEUS_URL,
		Path:       RULESPATH,
	}

	data, err := query.GetAlertManagerResponse()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
