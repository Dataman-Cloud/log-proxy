package api

import (
	"container/list"
	"errors"
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
	// DeleteAlertError delete keywrod error
	DeleteAlertError = "503-11005"
	// GetAlertError get keyword error
	GetAlertError = "503-11006"
	// UpdateAlertError update keywrod error
	UpdateAlertError = "503-11007"
	// GetEventsError get event history error
	GetEventsError = "503-11008"
	// GetPrometheusError get prometheus event error
	GetPrometheusError = "503-11009"
	// GetLogError get log error
	GetLogError = "503-11010"
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

// Applications get all applications
func (s *Search) Applications(ctx *gin.Context) {
	apps, err := s.Service.Applications(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAppsError, err))
		return
	}
	utils.Ok(ctx, apps)
}

// Tasks search applications tasks
func (s *Search) Tasks(ctx *gin.Context) {

	tasks, err := s.Service.Tasks(ctx.Param("app"), ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, tasks)
}

// Paths search applications paths
func (s *Search) Paths(ctx *gin.Context) {
	paths, err := s.Service.Paths(
		ctx.Query("cluster"),
		ctx.Query("user"),
		ctx.Param("app"),
		ctx.Query("task"),
		ctx.MustGet("page").(models.Page),
	)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, paths)
}

// Index search log by condition
func (s *Search) Index(ctx *gin.Context) {
	if ctx.Query("app") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("app can't be empty")))
		return
	}

	results, err := s.Service.Search(
		ctx.Query("cluster"),
		ctx.Query("user"),
		ctx.Query("app"),
		ctx.Query("task"),
		ctx.Query("path"),
		ctx.Query("keyword"),
		ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
}

// Context search log context
func (s *Search) Context(ctx *gin.Context) {
	if ctx.Query("app") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("app can't be empty")))
		return
	}

	if ctx.Query("task") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("task can't be empty")))
		return
	}

	if ctx.Query("path") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("path can't be empty")))
		return
	}

	if ctx.Query("offset") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("offset can't be empty")))
		return
	}

	results, err := s.Service.Context(
		ctx.Query("cluster"),
		ctx.Query("user"),
		ctx.Query("app"),
		ctx.Query("task"),
		ctx.Query("path"),
		ctx.Query("offset"),
		ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
}
