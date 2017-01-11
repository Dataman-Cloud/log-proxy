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

// Query return the result of metric query or expr query
func (m *Monitor) Query(ctx *gin.Context) {
	param := &backends.QueryParameter{
		Metric:    ctx.Query("metric"),
		ClusterID: ctx.Query("clusterid"),
		AppID:     ctx.Query("appid"),
		SlotID:    ctx.Query("taskid"), //SlotID is the swan's application field.
		UserID:    ctx.Query("userid"),
		Start:     ctx.Query("start"),
		End:       ctx.Query("end"),
		Step:      ctx.Query("step"),
		Period:    ctx.Query("period"),
		Expr:      ctx.Query("expr"),
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

	if param.Metric != "" && param.ClusterID == "" {
		err := fmt.Errorf("The paramter metric and clusterid required")
		utils.ErrorResponse(ctx, err)
		return
	}

	if param.Metric != "" && param.ClusterID != "" && param.UserID == "" {
		err := fmt.Errorf("The paramter userid required")
		utils.ErrorResponse(ctx, err)
		return
	}

	var err error
	param.SlotID, err = utils.ParseMonitorTaskID(param.SlotID)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	if param.Expr != "" {
		query := &backends.Query{
			HTTPClient:     http.DefaultClient,
			PromServer:     config.GetConfig().PrometheusURL,
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
			HTTPClient:     http.DefaultClient,
			PromServer:     config.GetConfig().PrometheusURL,
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

// QueryInfo return the info of clusters/cluster/app/node info
func (m *Monitor) QueryInfo(ctx *gin.Context) {
	param := &backends.QueryParameter{
		ClusterID: ctx.Query("clusterid"),
		UserID:    ctx.Query("userid"),
		AppID:     ctx.Query("appid"),
	}
	query := &backends.Query{
		HTTPClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PrometheusURL,
		Path:           backends.QUERYPATH,
		QueryParameter: param,
	}

	data := service.NewInfo()
	err := data.GetQueryInfo(query)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// QueryNodes return the metric data of nodes
func (m *Monitor) QueryNodes(ctx *gin.Context) {
	param := &backends.QueryParameter{
		ClusterID: ctx.Query("clusterid"),
		NodeID:    ctx.Query("nodeid"),
	}
	query := &backends.Query{
		HTTPClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PrometheusURL,
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
}

// GetAlerts return the Alerts
func (m *Monitor) GetAlerts(ctx *gin.Context) {
	query := &backends.AlertManager{
		HTTPClient: http.DefaultClient,
		Server:     config.GetConfig().AlertManagerURL,
		Path:       backends.AlertsPath,
	}
	data, err := query.GetAlertManagerResponse()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// GetAlertsGroups return the Alerts Groups queryed from Alertmanager
func (m *Monitor) GetAlertsGroups(ctx *gin.Context) {
	query := &backends.AlertManager{
		HTTPClient: http.DefaultClient,
		Server:     config.GetConfig().AlertManagerURL,
		Path:       backends.AlertsGroupsPath,
	}

	data, err := query.GetAlertManagerResponse()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// GetAlertsStatus return the Alerts Status queryed from Alertmanager
func (m *Monitor) GetAlertsStatus(ctx *gin.Context) {
	query := &backends.AlertManager{
		HTTPClient: http.DefaultClient,
		Server:     config.GetConfig().AlertManagerURL,
		Path:       backends.AlertsStatusPath,
	}

	data, err := query.GetAlertManagerResponse()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}
