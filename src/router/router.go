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
	logv1 := r.Group("/v1/search")
	{
		logv1.GET("/ping", s.Ping)
		logv1.GET("/applications", s.Applications)
		logv1.GET("/tasks/:app", s.Tasks)
		logv1.GET("/paths/:app", s.Paths)
		logv1.GET("/index", s.Index)
		logv1.GET("/context", s.Context)

		logv1.GET("/keyword", s.GetAlerts)
		logv1.POST("/keyword", s.CreateAlert)
		logv1.PUT("/keyword", s.UpdateAlert)
		logv1.DELETE("/keyword/:id", s.DeleteAlert)
		logv1.GET("/keyword/:id", s.GetAlert)

	}

	pv1 := r.Group("/v1/receive")
	{
		pv1.POST("/log", s.ReceiverLog)
	}

	monitor := api.NewMonitor()
	monitorv1 := r.Group("/v1/monitor")
	{
		monitorv1.GET("/query/items", monitor.GetQueryItems)
		monitorv1.GET("/users", monitor.GetUsers)
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
