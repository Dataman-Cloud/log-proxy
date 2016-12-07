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

func Router(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(gin.Recovery())
	r.Use(utils.Ginrus(log.StandardLogger(), time.RFC3339Nano, false))
	r.Use(middleware.CORSMiddleware())
	r.Use(middlewares...)

	s := api.GetSearch()
	logv1 := r.Group("/v1/search")
	{
		logv1.GET("/ping", s.Ping)
		logv1.GET("/applications", s.Applications)
		logv1.GET("/tasks/:appid", s.Tasks)
		logv1.GET("/paths/:appid", s.Paths)
		logv1.GET("/index", s.Index)
		logv1.GET("/context", s.Context)

	}

	av1 := r.Group("/v1/monitor")
	{
		av1.GET("/alert", s.GetAlerts)
		av1.POST("/alert", s.CreateAlert)
		av1.PUT("/alert", s.UpdateAlert)
		av1.DELETE("/alert/:id", s.DeleteAlert)
		av1.GET("/alert/:id", s.GetAlert)
		av1.GET("/prometheus", s.GetPrometheus)
	}

	pv1 := r.Group("/v1/recive")
	{
		pv1.POST("/prometheus", s.Receiver)
		pv1.POST("/log", s.ReceiverLog)
	}

	monitor := api.NewMonitor()
	monitorv1 := r.Group("/v1/monitor")
	{
		monitorv1.GET("/ping", monitor.Ping)
		// Query metric/expr
		monitorv1.GET("/query", monitor.Query)
		// Query info
		monitorv1.GET("/info", monitor.QueryInfo)
		monitorv1.GET("/nodes", monitor.QueryNodes)

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

	staticRouter := r.Group("/ui")
	{
		staticRouter.GET("/*filepath", func(ctx *gin.Context) {
			ctx.File("frontend/index.html")
		})
	}

	return r
}
