package router

import (
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/api"
	exporterapi "github.com/Dataman-Cloud/log-proxy/src/exporter/api"
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

	ep := exporterapi.NewExporter()

	exporterv1 := r.Group("/v1/exporter")
	{
		exporterv1.GET("/keyword", ep.GetFilters)
		exporterv1.POST("/keyword", ep.CreateFilter)
		exporterv1.PUT("/keyword", ep.UpdateFilter)
		exporterv1.DELETE("/keyword/:id", ep.DeleteFilter)
		exporterv1.GET("/keyword/:id", ep.GetFilter)
		exporterv1.POST("/input", ep.ReceiveLog)
	}

	s := api.NewSearch()
	logv1 := r.Group("/v1/search")
	{
		logv1.GET("/ping", s.Ping)
		logv1.GET("/applications", s.Applications)
		logv1.GET("/tasks/:app", s.Tasks)
		logv1.GET("/paths/:app", s.Paths)
		logv1.GET("/index", s.Index)
		logv1.GET("/context", s.Context)

		logv1.GET("/prometheus", s.GetPrometheus)
		logv1.GET("/prometheus/:id", s.GetPrometheu)

	}

	pv1 := r.Group("/v1/receive")
	{
		pv1.POST("/prometheus", s.Receiver)
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
