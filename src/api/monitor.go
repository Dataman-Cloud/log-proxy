package api

import (
	"fmt"
	"net/http"

	"github.com/Dataman-Cloud/log-proxy/src/backends"
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

// Ping return string "success"
func (m *Monitor) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

// Query return the result of metric query or expr query
/*
func (m *Monitor) Query(ctx *gin.Context) {
	param := &backends.QueryParameter{
		Metric:  ctx.Query("metric"),
		Cluster: ctx.Query("cluster"),
		App:     ctx.Query("app"),
		Task:    ctx.Query("task"), //Slot is the swan's application field.
		User:    ctx.Query("user"),
		Start:   ctx.Query("start"),
		End:     ctx.Query("end"),
		Step:    ctx.Query("step"),
		Period:  ctx.Query("period"),
		Expr:    ctx.Query("expr"),
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

	if param.Metric != "" && param.Cluster == "" {
		err := fmt.Errorf("The paramter metric and cluster required")
		utils.ErrorResponse(ctx, err)
		return
	}

	if param.Metric != "" && param.Cluster != "" && param.User == "" {
		err := fmt.Errorf("The paramter user required")
		utils.ErrorResponse(ctx, err)
		return
	}

	var err error
	param.Slot, err = utils.ParseMonitorTask(param.Slot)
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
*/

// GetQueryItems return the items of query metrics
func (m *Monitor) GetQueryItems(ctx *gin.Context) {
	query := &service.Query{
		ExprTmpl: service.SetQueryExprsList(),
	}
	utils.Ok(ctx, query.GetQueryItemList())
}

// GetClusters return the items of query metrics
func (m *Monitor) GetClusters(ctx *gin.Context) {
	param := &models.QueryParameter{
		Start: ctx.Query("start"),
		End:   ctx.Query("end"),
		Step:  ctx.Query("step"),
	}

	query := &service.Query{
		ExprTmpl:       service.SetQueryExprsList(),
		HTTPClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PrometheusURL,
		Path:           backends.QUERYRANGEPATH,
		QueryParameter: param,
	}

	data, err := query.GetQueryClusters()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// GetClusters return the items of query metrics
func (m *Monitor) GetClusterApps(ctx *gin.Context) {
	param := &models.QueryParameter{
		Start:   ctx.Query("start"),
		End:     ctx.Query("end"),
		Step:    ctx.Query("step"),
		Cluster: ctx.Param("clusterid"),
	}

	query := &service.Query{
		ExprTmpl:       service.SetQueryExprsList(),
		HTTPClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PrometheusURL,
		Path:           backends.QUERYRANGEPATH,
		QueryParameter: param,
	}

	data, err := query.GetQueryApps()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// GetClusters return the items of query metrics
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
		Path:           backends.QUERYRANGEPATH,
		QueryParameter: param,
	}

	data, err := query.GetQueryAppTasks()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

func (m *Monitor) Query(ctx *gin.Context) {
	param := &models.QueryParameter{
		Metric:  ctx.Query("metric"),
		Cluster: ctx.Query("cluster"),
		App:     ctx.Query("app"),
		Task:    ctx.Query("task"),
		User:    ctx.Query("user"),
		Start:   ctx.Query("start"),
		End:     ctx.Query("end"),
		Step:    ctx.Query("step"),
		Expr:    ctx.Query("expr"),
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

	if param.Metric != "" && param.Cluster == "" {
		err := fmt.Errorf("The paramter metric and cluster required")
		utils.ErrorResponse(ctx, err)
		return
	}
	/*
		if param.Metric != "" && param.Cluster != "" && param.User == "" {
			err := fmt.Errorf("The paramter user required")
			utils.ErrorResponse(ctx, err)
			return
		}
	*/
	if param.Expr != "" {
		query := &service.Query{
			ExprTmpl:       service.SetQueryExprsList(),
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
		utils.Ok(ctx, data)
		return
	}

	if param.Metric != "" {
		query := &service.Query{
			ExprTmpl:       service.SetQueryExprsList(),
			HTTPClient:     http.DefaultClient,
			PromServer:     config.GetConfig().PrometheusURL,
			Path:           backends.QUERYRANGEPATH,
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

// QueryInfo return the info of clusters/cluster/app/node info
func (m *Monitor) QueryInfo(ctx *gin.Context) {
	param := &backends.QueryParameter{
		Cluster: ctx.Query("cluster"),
		User:    ctx.Query("user"),
		App:     ctx.Query("app"),
		Slot:    ctx.Query("task"), //Slot is the swan's application field.
	}

	var err error
	param.Slot, err = utils.ParseMonitorTask(param.Slot)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	query := &backends.Query{
		HTTPClient:     http.DefaultClient,
		PromServer:     config.GetConfig().PrometheusURL,
		Path:           backends.QUERYPATH,
		QueryParameter: param,
	}

	data := service.NewInfo()
	err = data.GetQueryInfo(query)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// QueryNodes return the metric data of nodes
func (m *Monitor) QueryNodes(ctx *gin.Context) {
	param := &backends.QueryParameter{
		Cluster: ctx.Query("cluster"),
		Node:    ctx.Query("node"),
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
