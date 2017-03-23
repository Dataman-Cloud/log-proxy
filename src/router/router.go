package router

import (
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/api"
	"github.com/Dataman-Cloud/log-proxy/src/router/middleware"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router add router function
func Router(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(gin.Recovery())
	r.Use(utils.Ginrus(log.StandardLogger(), time.RFC3339Nano, false))
	r.Use(middleware.CORSMiddleware())
	r.Use(middlewares...)

	s := api.GetSearch()
	logv1 := r.Group("/v1/log")
	{
		logv1.GET("/ping", s.Ping)
		logv1.GET("/clusters", s.Clusters)
		logv1.GET("/clusters/:cluster/apps", s.Applications)
		logv1.GET("/clusters/:cluster/apps/:app/tasks", s.Tasks)
		logv1.GET("/paths/:app", s.Paths)
		logv1.GET("/index", s.Index)
		logv1.GET("/context", s.Context)

		logv1.GET("/rules", s.GetAlerts)
		logv1.POST("/rules", s.CreateAlert)
		logv1.PUT("/rules", s.UpdateAlert)
		logv1.DELETE("/rules/:id", s.DeleteAlert)
		logv1.GET("/rules/:id", s.GetAlert)

	}

	pv1 := r.Group("/v1/receive")
	{
		pv1.POST("/log", s.ReceiverLog)
	}

	monitor := api.NewMonitor()
	monitorv1 := r.Group("/v1/monitor")
	{
		monitorv1.GET("/ping", monitor.Ping)
		// Query metric/expr
		monitorv1.GET("/query", monitor.Query)
		monitorv1.GET("/query/metrics", monitor.GetQueryItems)
		monitorv1.GET("/cluster", monitor.GetClusters)
		monitorv1.GET("/cluster/:clusterid/apps", monitor.GetClusterApps)
		monitorv1.GET("/cluster/:clusterid/apps/:appid/tasks", monitor.GetAppsTasks)
		// Query info
		//monitorv1.GET("/info", monitor.QueryInfo)
		//monitorv1.GET("/nodes", monitor.QueryNodes)

		// AlertManager API
		monitorv1.GET("/alerts", monitor.GetAlerts)
		monitorv1.GET("/alerts/groups", monitor.GetAlertsGroups)
		monitorv1.GET("/alerts/status", monitor.GetAlertsStatus)
		monitorv1.GET("/silences", monitor.GetSilences)
		monitorv1.POST("/silences", monitor.CreateSilence)
		monitorv1.GET("/silence/:id", monitor.GetSilence)
		monitorv1.DELETE("/silence/:id", monitor.DeleteSilence)
		monitorv1.PUT("/silence/:id", monitor.UpdateSilence)
	}

	alert := api.NewAlert()
	alertv1 := r.Group("/v1/alert")
	{
		alertv1.POST("/rules", alert.CreateAlertRule)
		alertv1.DELETE("/rules/:id", alert.DeleteAlertRule)
		alertv1.GET("/rules", alert.ListAlertRules)
		alertv1.GET("/rules/:id", alert.GetAlertRule)
		alertv1.PUT("/rules", alert.UpdateAlertRule)
		alertv1.POST("/rules/conf", alert.ReloadAlertRuleConf)
		alertv1.POST("/receiver", alert.ReceiveAlertEvent)
		alertv1.PUT("/events/:id", alert.AckAlertEvent)
		alertv1.GET("/events", alert.GetAlertEvents)
	}
	alert.AlertRuleFilesMaintainer()

	staticRouter := r.Group("/ui")
	{
		staticRouter.GET("/*filepath", func(ctx *gin.Context) {
			ctx.File("frontend/index.html")
		})
	}

	return r
}
