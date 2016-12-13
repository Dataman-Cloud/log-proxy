package api

import (
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

const (
	GETAPPS_ERROR        = "503-11000"
	PARAM_ERROR          = "400-11001"
	GETTASK_ERROR        = "503-11002"
	INDEX_ERROR          = "503-11003"
	CREATE_ALERT_ERROR   = "503-11004"
	DELETE_ALERT_ERROR   = "503-11005"
	GET_ALERT_ERROR      = "503-11006"
	UPDATE_ALERT_ERROR   = "503-11007"
	GET_EVENTS_ERROR     = "503-11008"
	GET_PROMETHEUS_ERROR = "503-11009"
	GET_LOG_ERROR        = "503-11010"
)

type search struct {
	Service       *service.SearchService
	KeywordFilter map[string][]string
	Counter       *prometheus.CounterVec
	Kmutex        *sync.Mutex
}

func GetSearch() *search {
	s := &search{
		Service:       service.NewEsService(strings.Split(config.GetConfig().ES_URL, ",")),
		KeywordFilter: make(map[string][]string),
		Counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "log_keyword",
				Help: "log keyword counter",
			},
			[]string{"appid", "taskid", "path", "keyword", "userid", "clusterid", "offset"},
		),
		Kmutex: new(sync.Mutex),
	}
	prometheus.MustRegister(s.Counter)

	if s.Service == nil {
		return s
	}

	alerts, err := s.Service.GetAlerts(models.Page{
		PageFrom: 0,
		PageSize: 0,
	})

	if err != nil || alerts == nil {
		return s
	}

	alerts, err = s.Service.GetAlerts(models.Page{
		PageFrom: 0,
		PageSize: int(alerts["count"].(int64)),
	})

	if err != nil {
		return s
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	for _, alert := range alerts["results"].([]models.Alert) {
		s.KeywordFilter[alert.AppId+alert.Path] = append(s.KeywordFilter[alert.AppId+alert.Path], alert.Keyword)
	}

	return s
}

func (s *search) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

func (s *search) Applications(ctx *gin.Context) {
	apps, err := s.Service.Applications(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GETAPPS_ERROR, err))
		return
	}
	utils.Ok(ctx, apps)
}

func (s *search) Tasks(ctx *gin.Context) {
	appName := ctx.Param("appid")
	if appName == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("param error")))
		return
	}

	tasks, err := s.Service.Tasks(appName, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GETTASK_ERROR, err))
		return
	}
	utils.Ok(ctx, tasks)
}

func (s *search) Paths(ctx *gin.Context) {
	appName := ctx.Param("appid")
	if appName == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("param error")))
		return
	}

	paths, err := s.Service.Paths(appName, ctx.Query("taskid"), ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GETTASK_ERROR, err))
		return
	}
	utils.Ok(ctx, paths)
}

func (s *search) Index(ctx *gin.Context) {
	if ctx.Query("appid") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("appid can't be empty")))
		return
	}

	results, err := s.Service.Search(ctx.Query("appid"),
		ctx.Query("taskid"),
		ctx.Query("path"),
		ctx.Query("keyword"),
		ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(INDEX_ERROR, err))
		return
	}

	utils.Ok(ctx, results)
}

func (s *search) Context(ctx *gin.Context) {
	if ctx.Query("appid") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("appid can't be empty")))
		return
	}

	if ctx.Query("taskid") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("taskid can't be empty")))
		return
	}

	if ctx.Query("path") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("path can't be empty")))
		return
	}

	if ctx.Query("offset") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("offset can't be empty")))
		return
	}

	results, err := s.Service.Context(ctx.Query("appid"),
		ctx.Query("taskid"),
		ctx.Query("path"),
		ctx.Query("offset"),
		ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(INDEX_ERROR, err))
		return
	}

	utils.Ok(ctx, results)
}
