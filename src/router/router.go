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
	r.GET("/ping", s.Ping)

	logRouter := r.Group("/v2/log")
	{
		logRouter.GET("/clusters", s.Clusters)
		logRouter.GET("/clusters/:cluster/apps", s.Applications)
		logRouter.GET("/clusters/:cluster/apps/:app/slots", s.Slots)
		logRouter.GET("/clusters/:cluster/apps/:app/slots/:slot/tasks", s.Tasks)
		logRouter.GET("/clusters/:cluster/apps/:app/sources", s.Sources)
		logRouter.GET("/clusters/:cluster/apps/:app/search", s.Search)
		logRouter.GET("/clusters/:cluster/apps/:app/context", s.Context)

		logRouter.POST("/alert/rules", s.CreateLogAlertRule)
		//logRouter.GET("/alert/rules", s.GetLogAlertRules)
		//logRouter.GET("/alert/rules/:id", s.GetLogAlertRule)
		logRouter.PUT("/alert/rules/:id", s.UpdateLogAlertRule)
		//logRouter.DELETE("/alert/rules/:id", s.DeleteLogAlertRule)
	}

	pv1 := r.Group("/v1/receive")
	{
		pv1.POST("/log", s.ReceiverLog)
	}

	monitor := api.NewMonitor()
	monitorv1 := r.Group("/v1/monitor")
	{
		monitorv1.GET("/query", monitor.Query)
		monitorv1.GET("/query/items", monitor.GetQueryItems)
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
