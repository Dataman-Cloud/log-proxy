package api

import (
	"container/list"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

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
	// GetLogError get log error
	GetLogError = "503-11010"
	// ParamError param error
	ParamError = "400-11001"
	// CreateFilterError create keyword error
	CreateFilterError = "503-11004"
	// DeleteFilterError delete keyword error
	DeleteFilterError = "503-11005"
	// GetFilterError get keyword error
	GetFilterError = "503-11006"
	// UpdateFilterError update keyword error
	UpdateFilterError = "503-11007"
)

// Exporter prometheus exporter client struct
type Exporter struct {
	ES            *service.SearchService
	KeywordFilter map[string]*list.List
	Counter       *prometheus.CounterVec
	Kmutex        *sync.Mutex
}

// NewExporter new search client
func NewExporter() *Exporter {
	ep := &Exporter{
		ES:            service.NewEsService(strings.Split(config.GetConfig().EsURL, ",")),
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
		prometheus.MustRegister(ep.Counter)
		registry = true
	}

	filters, err := ep.ES.GetFilters(models.Page{
		PageFrom: 0,
		PageSize: 1000,
	})

	if err != nil {
		return ep
	}

	ep.Kmutex.Lock()
	defer ep.Kmutex.Unlock()
	if filters == nil {
		return ep
	}
	for _, filter := range filters["results"].([]models.KWFilter) {
		if ep.KeywordFilter[filter.AppID+filter.Path] == nil {
			ep.KeywordFilter[filter.AppID+filter.Path] = list.New()
		}
		ep.KeywordFilter[filter.AppID+filter.Path].PushBack(filter.Keyword)
	}

	return ep
}

// CreateFilter create keyword filter
func (ep *Exporter) CreateFilter(ctx *gin.Context) {
	filter := new(models.KWFilter)
	if err := ctx.BindJSON(filter); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("request body param error")))
		return
	}

	if filter.AppID == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("appid can't be empty")))
		return
	}

	if filter.Keyword == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("keyword can't be empty")))
		return
	}

	if ep.KeywordFilter[filter.AppID+filter.Path] != nil {
		for e := ep.KeywordFilter[filter.AppID+filter.Path].Front(); e != nil; e = e.Next() {
			if e.Value.(string) == filter.Keyword {
				utils.ErrorResponse(ctx, utils.NewError(CreateFilterError, errors.New("keyword exist")))
				return
			}
		}
	}

	err := ep.ES.CreateFilter(filter)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(CreateFilterError, err))
		return
	}

	ep.Kmutex.Lock()
	defer ep.Kmutex.Unlock()
	ep.KeywordFilter[filter.AppID+filter.Path].PushBack(filter.Keyword)

	utils.Ok(ctx, "create success")
}

// DeleteFilter delete keyword filter
func (ep *Exporter) DeleteFilter(ctx *gin.Context) {

	filter, err := ep.ES.GetFilter(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetFilterError, err))
		return
	}

	err = ep.ES.DeleteFilter(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(DeleteFilterError, err))
		return
	}

	ep.Kmutex.Lock()
	defer ep.Kmutex.Unlock()
	for e := ep.KeywordFilter[filter.AppID+filter.Path].Front(); e != nil; e = e.Next() {
		if e.Value.(string) == filter.Keyword {
			ep.KeywordFilter[filter.AppID+filter.Path].Remove(e)
		}
	}

	utils.Ok(ctx, "delete success")
}

// GetFilters get all keyword filter
func (ep *Exporter) GetFilters(ctx *gin.Context) {
	results, err := ep.ES.GetFilters(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetFilterError, err))
		return
	}

	utils.Ok(ctx, results)
}

// GetFilter get keyword filter by id
func (ep *Exporter) GetFilter(ctx *gin.Context) {
	result, err := ep.ES.GetFilter(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetFilterError, err))
		return
	}

	utils.Ok(ctx, result)
}

// UpdateFilter update keyword filter
func (ep *Exporter) UpdateFilter(ctx *gin.Context) {
	filter := new(models.KWFilter)
	if err := ctx.BindJSON(filter); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("request body param error")))
		return
	}

	if filter.ID == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("id can't be empty")))
		return
	}

	if filter.AppID == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("appid can't be empty")))
		return
	}

	if filter.Keyword == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("keyword can't be empty")))
		return
	}

	result, err := ep.ES.GetFilter(filter.ID)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetFilterError, err))
		return
	}

	filter.CreateTime = time.Now().Format(time.RFC3339Nano)
	err = ep.ES.UpdateFilter(filter)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(UpdateFilterError, err))
		return
	}

	ep.Kmutex.Lock()
	defer ep.Kmutex.Unlock()
	for e := ep.KeywordFilter[result.AppID+result.Path].Front(); e != nil; e = e.Next() {
		if e.Value.(string) == filter.Keyword {
			ep.KeywordFilter[result.AppID+result.Path].Remove(e)
			ep.KeywordFilter[result.AppID+result.Path].PushBack(filter.Keyword)
			break
		}
	}

	utils.Ok(ctx, "update success")
}

// ReceiveLog receive log data from logstash
func (ep *Exporter) ReceiveLog(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	app, ok := m["app"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found app"))
		return
	}

	task, ok := m["task"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found task"))
		return
	}

	path, ok := m["path"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found path"))
		return
	}

	user, ok := m["user"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found user"))
		return
	}

	cluster, ok := m["cluster"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found cluster"))
		return
	}

	message, ok := m["message"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found message"))
		return
	}

	keywords, ok := ep.KeywordFilter[app.(string)+path.(string)]
	if !ok {
		utils.Ok(ctx, "ok")
		return
	}
	for e := keywords.Front(); e != nil; e = e.Next() {
		if strings.Index(message.(string), e.Value.(string)) == -1 {
			continue
		}

		ep.Counter.WithLabelValues(
			app.(string),
			task.(string),
			path.(string),
			e.Value.(string),
			user.(string),
			cluster.(string),
		).Inc()
	}

	utils.Ok(ctx, "ok")
	return
}
