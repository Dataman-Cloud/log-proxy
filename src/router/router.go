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
	go s.PollAlert()
	logv1 := r.Group("/v1/search")
	{
		logv1.GET("/ping", s.Ping)
		logv1.GET("/applications", s.Applications)
		logv1.GET("/tasks/:appid", s.Tasks)
		logv1.GET("/paths/:appid/:taskid", s.Paths)
		logv1.GET("/index", s.Index)
		logv1.GET("/context", s.Context)

		logv1.GET("/alert", s.GetAlerts)
		logv1.POST("/alert", s.CreateAlert)
		logv1.PUT("/alert", s.UpdateAlert)
		logv1.DELETE("/alert/:id", s.DeleteAlert)
		logv1.POST("/alert/prometheus", s.Receiver)
		logv1.GET("/alert/keyword/hostory", s.GetKeywordAlertHistory)
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
	}

	staticRouter := r.Group("/ui")
	{
		staticRouter.GET("/*filepath", func(ctx *gin.Context) {
			ctx.File("frontend/index.html")
		})
	}

	return r
}
