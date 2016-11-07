package router

import (
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/api"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func Router(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(utils.Ginrus(log.StandardLogger(), time.RFC3339Nano, false))
	r.Use(middlewares...)

	s := api.GetSearch()
	logv1 := r.Group("/v1/search")
	{
		logv3.GET("/ping", s.Ping)
	}

	monitor := api.GetMonitor()
	monitorv1 := r.Group("/v1/monitor")
	{
		monitorv3.GET("/ping", monitor.Ping)
	}

	return r
}
