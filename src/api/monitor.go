package api

import (
	"fmt"
	"net/http"

	"github.com/Dataman-Cloud/log-proxy/src/backends"
	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

// Backends API: AlertManager
const (
	ALERTSPATH       = "/api/v1/alerts"
	ALERTSGROUSPPATH = "/api/v1/alerts/groups"
	ALERTSSTATUSPATH = "/api/v1/status"
)

type monitor struct {
}

func NewMonitor() *monitor {
	return &monitor{}
}

func (m *monitor) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

func (m *monitor) Query(ctx *gin.Context) {
	param := &backends.QueryParameter{
		Metric:    ctx.Query("metric"),
		ClusterID: ctx.Query("clusterid"),
		AppID:     ctx.Query("appid"),
		TaskID:    ctx.Query("taskid"),
		Start:     ctx.Query("start"),
		End:       ctx.Query("end"),
		Step:      ctx.Query("step"),
		Period:    ctx.Query("period"),
		Expr:      ctx.Query("expr"),
	}

	if param.Metric != "" && param.Expr != "" {
		err := fmt.Errorf("The paramter confict between metric and expr!")
		utils.ErrorResponse(ctx, err)
		return
	}

	if param.Metric == "" && param.Expr == "" {
		err := fmt.Errorf("The paramter of metric or expr required!")
		utils.ErrorResponse(ctx, err)
		return
	}

	if param.Expr != "" {
		query := &backends.Query{
			HttpClient:     http.DefaultClient,
			PromServer:     config.GetConfig().PROMETHEUS_URL,
			Path:           backends.QUERYRANGEPATH,
			QueryParameter: param,
		}

		data, err := query.QueryExpr()
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, data)
	}

	if param.Metric != "" {
		query := &backends.Query{
			HttpClient:     http.DefaultClient,
			PromServer:     config.GetConfig().PROMETHEUS_URL,
			Path:           backends.QUERYRANGEPATH,
			QueryParameter: param,
		}
		data := service.NewMetric()
		err := data.GetQueryMetric(query)
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
		utils.Ok(ctx, data)
	}
}

func (m *monitor) QueryInfo(ctx *gin.Context) {
	param := &backends.QueryParameter{
		ClusterID: ctx.Query("clusterid"),
		AppID:     ctx.Query("appid"),
	}
	query := &backends.Query{
		HttpClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PROMETHEUS_URL,
		Path:           backends.QUERYPATH,
		QueryParameter: param,
	}

	if query.ClusterID != "" && query.AppID != "" {
		err := fmt.Errorf("The paramter confict between clusterid and appid!")
		utils.ErrorResponse(ctx, err)
		return
	}

	data := service.NewInfo()
	err := data.GetQueryInfo(query)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

func (m *monitor) QueryNodes(ctx *gin.Context) {
	param := &backends.QueryParameter{
		ClusterID: ctx.Query("clusterid"),
		NodeID:    ctx.Query("nodeid"),
	}
	query := &backends.Query{
		HttpClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PROMETHEUS_URL,
		Path:           backends.QUERYRANGEPATH,
		QueryParameter: param,
	}

	data := service.NewNodesInfo()
	err := data.GetQueryNodesInfo(query)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)

	return
}

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
	utils.Ok(ctx, data)
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
	utils.Ok(ctx, data)
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
	utils.Ok(ctx, data)
}
