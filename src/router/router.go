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
	s.InitLogKeywordFilter()

	r.GET("/ping", s.Ping)

	logRouter := r.Group("/v2/log")
	{
		logRouter.GET("/apps", s.Applications)
		logRouter.GET("/apps/:app/slots", s.Slots)
		logRouter.GET("/apps/:app/slots/:slot/tasks", s.Tasks)
		logRouter.GET("/apps/:app/sources", s.Sources)
		logRouter.GET("/apps/:app/search", s.Search)
		logRouter.GET("/apps/:app/context", s.Context)

		logRouter.GET("/everything/:key", s.Everything)

		logRouter.POST("/alert/rules", s.CreateLogAlertRule)
		logRouter.GET("/alert/rules", s.GetLogAlertRules)
		logRouter.GET("/alert/rules/:id", s.GetLogAlertRule)
		logRouter.PUT("/alert/rules/:id", s.UpdateLogAlertRule)
		logRouter.DELETE("/alert/rules/:id", s.DeleteLogAlertRule)
		logRouter.GET("/alert/events", s.GetLogAlertEvents)
		logRouter.GET("/alert/events/apps", s.GetLogAlertApps)
		logRouter.PATCH("/alert/events/:id", s.LogAlertEventAction)
		logRouter.POST("/receive", s.ReceiveLog)
	}

	monitor := api.NewMonitor()
	monitorv1 := r.Group("/v1/monitor")
	{
		monitorv1.GET("/query", monitor.Query)
		monitorv1.GET("/query/items", monitor.GetQueryItems)
		monitorv1.GET("/apps", monitor.GetApps)
		monitorv1.GET("/apps/:appid/tasks", monitor.GetAppsTasks)
	}

	alertv1 := r.Group("/v1/alert")
	{
		alertv1.GET("/indicators", monitor.GetAlertIndicators)
		alertv1.POST("/rules", monitor.CreateAlertRule)
		alertv1.DELETE("/rules/:id", monitor.DeleteAlertRule)
		alertv1.GET("/rules", monitor.ListAlertRules)
		alertv1.GET("/rules/:id", monitor.GetAlertRule)
		alertv1.PUT("/rules/:id", monitor.UpdateAlertRule)

		alertv1.POST("/receiver", monitor.ReceiveAlertEvent)
		alertv1.GET("/events", monitor.GetAlertEvents)
		alertv1.PUT("/events/:id", monitor.AckAlertEvent)
		/*
			alertv1.POST("/rules/conf", alert.ReloadAlertRuleConf)
			alertv1.POST("/receiver", alert.ReceiveAlertEvent)
			alertv1.PUT("/events/:id", alert.AckAlertEvent)
			alertv1.GET("/events", alert.GetAlertEvents)
		*/
	}
	//alert.AlertRuleFilesMaintainer()

	staticRouter := r.Group("/ui")
	{
		staticRouter.GET("/*filepath", func(ctx *gin.Context) {
			ctx.File("frontend/index.html")
		})
	}

	return r
}
