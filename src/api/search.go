package api

import (
	"container/list"
	"strings"
	"sync"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// registry prometheus registry counter
var registry bool

const (
	// GetAppsError get application error
	GetAppsError = "503-11000"
	// ParamError param error
	ParamError = "400-11001"
	// GetTaskError get task error
	GetTaskError = "503-11002"
	// IndexError search index log error
	IndexError = "503-11003"
	// CreateAlertError create keyword error
	CreateAlertError = "503-11004"
	// DeleteAlertError delete keyword error
	DeleteAlertError = "503-11005"
	// GetAlertError get keyword error
	GetAlertError = "503-11006"
	// UpdateAlertError update keyword error
	UpdateAlertError = "503-11007"
	// GetEventsError get event history error
	GetEventsError = "503-11008"
	// GetPrometheusError get prometheus event error
	GetPrometheusError = "503-11009"
	// GetLogError get log error
	GetLogError = "503-11010"

	GetClustersError = "503-11011"

	GetSlotsError = "503-11012"
)

// Search search client struct
type Search struct {
	Service       *service.SearchService
	KeywordFilter map[string]*list.List
	Counter       *prometheus.CounterVec
	Kmutex        *sync.Mutex
}

// GetSearch new search client
func GetSearch() *Search {
	s := &Search{
		Service:       service.NewEsService(strings.Split(config.GetConfig().EsURL, ",")),
		KeywordFilter: make(map[string]*list.List),
		Counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "log_keyword",
				Help: "log keyword counter",
			},
			[]string{"app", "task", "path", "keyword", "user", "cluster"},
		),
		Kmutex: new(sync.Mutex),
	}

	if !registry {
		prometheus.MustRegister(s.Counter)
		registry = true
	}

	alerts, err := s.Service.GetAlerts(models.Page{
		PageFrom: 0,
		PageSize: 1000,
	})

	if err != nil {
		return s
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	if alerts == nil {
		return s
	}
	for _, alert := range alerts["results"].([]models.Alert) {
		if s.KeywordFilter[alert.AppID+alert.Path] == nil {
			s.KeywordFilter[alert.AppID+alert.Path] = list.New()
		}
		s.KeywordFilter[alert.AppID+alert.Path].PushBack(alert.Keyword)
	}

	return s
}

// Ping ping
func (s *Search) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

func (s *Search) Clusters(ctx *gin.Context) {
	clusters, err := s.Service.Clusters(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetClustersError, err))
		return
	}

	utils.Ok(ctx, clusters)
	return
}

// Applications get all applications
func (s *Search) Applications(ctx *gin.Context) {
	cluster := ctx.Param("cluster")
	apps, err := s.Service.Applications(cluster, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAppsError, err))
		return
	}
	utils.Ok(ctx, apps)
}

func (s *Search) Slots(ctx *gin.Context) {
	cluster := ctx.Param("cluster")
	app := ctx.Param("app")

	slots, err := s.Service.Slots(cluster, app, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetSlotsError, err))
		return
	}

	utils.Ok(ctx, slots)
	return
}

// Tasks search applications tasks
func (s *Search) Tasks(ctx *gin.Context) {
	cluster := ctx.Param("cluster")
	app := ctx.Param("app")
	slot := ctx.Param("slot")
	tasks, err := s.Service.Tasks(cluster, app, slot, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, tasks)
}

// Paths search applications paths
func (s *Search) Sources(ctx *gin.Context) {
	cluster := ctx.Param("cluster")
	app := ctx.Param("app")

	options := make(map[string]interface{})
	options["page"] = ctx.MustGet("page")
	if ctx.Query("slot") != "" {
		options["slot"] = ctx.Query("slot")
	}

	if ctx.Query("task") != "" {
		options["task"] = ctx.Query("task")
	}

	sources, err := s.Service.Sources(cluster, app, options)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, sources)
}

func (s *Search) Search(ctx *gin.Context) {
	cluster := ctx.Param("cluster")
	app := ctx.Param("app")

	options := make(map[string]interface{})
	options["page"] = ctx.MustGet("page")
	if ctx.Query("slot") != "" {
		options["slot"] = ctx.Query("slot")
	}

	if ctx.Query("task") != "" {
		options["task"] = ctx.Query("task")
	}

	if ctx.Query("source") != "" {
		options["source"] = ctx.Query("source")
	}

	if ctx.Query("keyword") != "" {
		options["keyword"] = ctx.Query("keyword")
	}

	if ctx.Query("conj") != "" {
		options["conj"] = ctx.Query("conj")
	}

	results, err := s.Service.Search(cluster, app, options)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
}

// Context search log context
func (s *Search) Context(ctx *gin.Context) {
	cluster := ctx.Param("cluster")
	app := ctx.Param("app")

	options := make(map[string]interface{})
	options["page"] = ctx.MustGet("page")
	if ctx.Query("slot") != "" {
		options["slot"] = ctx.Query("slot")
	}

	if ctx.Query("task") != "" {
		options["task"] = ctx.Query("task")
	}

	if ctx.Query("source") != "" {
		options["source"] = ctx.Query("source")
	}

	if ctx.Query("offset") != "" {
		options["offset"] = ctx.Query("offset")
	}

	results, err := s.Service.Context(cluster, app, options)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
}
