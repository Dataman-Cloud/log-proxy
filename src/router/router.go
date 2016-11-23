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
	logv1 := r.Group("/v1/search", s.Middleware)
	{
		logv1.GET("/ping", s.Ping)
		logv1.GET("/applications", s.Applications)
		logv1.GET("/tasks/:appid", s.Tasks)
		logv1.GET("/paths/:appid/:taskid", s.Paths)
		logv1.GET("/index", s.Index)

		logv1.GET("/alert", s.GetAlerts)
		logv1.POST("/alert", s.CreateAlert)
		logv1.PUT("/alert", s.UpdateAlert)
		logv1.DELETE("/alert/:id", s.DeleteAlert)
	}

	monitor := api.NewMonitor()
	monitorv1 := r.Group("/v1/monitor")
	{
		monitorv1.GET("/ping", monitor.Ping)
		monitorv1.GET("", monitor.QueryRange)
		monitorv1.GET("/applications", monitor.QueryApps)
		monitorv1.GET("/nodes", monitor.QueryNodes)
		monitorv1.GET("/application", monitor.QueryApp)

		monitorv1.GET("/alerts", monitor.GetAlerts)
		monitorv1.GET("/alerts/groups", monitor.GetAlertsGroups)
		monitorv1.GET("/alerts/status", monitor.GetAlertsStatus)
	}

	return r
}
