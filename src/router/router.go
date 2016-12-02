package router

import (
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/api"
	"github.com/Dataman-Cloud/log-proxy/src/router/middleware"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func Router(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(utils.Ginrus(log.StandardLogger(), time.RFC3339Nano, false))
	r.Use(middleware.CORSMiddleware())
	r.Use(middlewares...)

	s := api.GetSearch()
	go s.ReceiverMarathonEvent()
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
		av1.GET("/events", s.GetEvents)
		av1.GET("/mointor", s.GetPrometheus)
	}

	pv1 := r.Group("/v1/recive")
	{
		pv1.POST("/prometheus", s.Receiver)
	}

	monitor := api.NewMonitor()
	monitorv1 := r.Group("/v1/monitor")
	{
		monitorv1.GET("/ping", monitor.Ping)
		monitorv1.GET("", monitor.QueryRange)
		monitorv1.GET("/applications", monitor.QueryApps)
		monitorv1.GET("/nodes", monitor.QueryNodes)
		monitorv1.GET("/application", monitor.QueryApp)

		// Promethues HTTP API
		monitorv1.GET("/promql/query", monitor.PromqlQuery)
		monitorv1.GET("/promql/query_range", monitor.PromqlQueryRange)
		// AlertManager API
		monitorv1.GET("/alerts", monitor.GetAlerts)
		monitorv1.GET("/alerts/groups", monitor.GetAlertsGroups)
		monitorv1.GET("/alerts/status", monitor.GetAlertsStatus)
		// Rules
		monitorv1.GET("/alerts/rules", monitor.GetAlertsRules)
	}

	staticRouter := r.Group("/ui")
	{
		staticRouter.GET("/*filepath", func(ctx *gin.Context) {
			ctx.File("frontend/index.html")
		})
	}

	return r
}
