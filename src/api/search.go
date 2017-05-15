package api

import (
	"container/list"
	"strings"
	"sync"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
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

	CreateLogAlertRuleError = "503-11004"

	DeleteLogAlertRuleError = "503-11005"

	GetLogAlertRuleError = "503-11006"

	UpdateLogAlertRuleError = "503-11007"

	GetEventsError = "503-11008"
	// GetPrometheusError get prometheus event error
	GetPrometheusError = "503-11009"
	// GetLogError get log error
	GetLogError = "503-11010"

	GetSlotsError      = "503-11011"
	GetLogContextError = "503-11012"

	GetLogAlertEventsError = "503-11013"
	GetLogAlertAppsError   = "503-11014"
)

// Search search client struct
type Search struct {
	Service       service.LogSearchService
	Store         store.Store
	KeywordFilter map[string]*list.List
	Kmutex        sync.RWMutex
}

// GetSearch new search client
func GetSearch() *Search {
	return &Search{
		Service:       service.NewEsService(strings.Split(config.GetConfig().EsURL, ",")),
		Store:         datastore.From(database.GetDB()),
		KeywordFilter: make(map[string]*list.List),
	}
}

func (s *Search) InitLogKeywordFilter() {
	opts := map[string]interface{}{}
	rules, err := s.Store.GetLogAlertRules(opts, models.Page{PageFrom: 0, PageSize: 10000})
	if err != nil {
		log.Errorf("get log alert ruels forn store failed. Error: %+v", err)
		return
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()

	if rules == nil {
		return
	}

	for _, rule := range rules["rules"].([]*models.LogAlertRule) {
		ruleIndex := getLogAlertRuleIndex(*rule)
		if s.KeywordFilter[ruleIndex] == nil {
			s.KeywordFilter[ruleIndex] = list.New()
		}
		s.KeywordFilter[ruleIndex].PushBack(*rule)
	}

	return
}

func getLogAlertRuleIndex(r models.LogAlertRule) string {
	return r.Group + "-" + r.App + "-" + r.Source
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

func (s *Search) Slots(ctx *gin.Context) {
	app := ctx.Param("app")
	slots, err := s.Service.Slots(app, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetSlotsError, err))
		return
	}

	utils.Ok(ctx, slots)
	return
}

// Tasks search applications tasks
func (s *Search) Tasks(ctx *gin.Context) {
	app := ctx.Param("app")
	slot := ctx.Param("slot")
	appLabel := config.LogAppLabel()
	slotLabel := config.LogSlotLabel()
	options := config.ConvertRequestQueryParams(ctx.Request.URL.Query())
	options[appLabel] = app
	options[slotLabel] = slot
	tasks, err := s.Service.Tasks(options, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, tasks)
}

// Paths search applications paths
func (s *Search) Sources(ctx *gin.Context) {
	app := ctx.Param("app")
	appLabel := config.LogAppLabel()
	options := config.ConvertRequestQueryParams(ctx.Request.URL.Query())
	options[appLabel] = app
	sources, err := s.Service.Sources(options, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, sources)
}

func (s *Search) Search(ctx *gin.Context) {
	keyword := ctx.Query(config.LogKeywordLabel())
	if keyword != "" {
		if strings.ToLower(ctx.Query(config.LogConjLabel())) == "or" {
			keyword = strings.Join(strings.Split(keyword, " "), " OR ")
		} else {
			keyword = strings.Join(strings.Split(keyword, " "), " AND ")
		}
	}
	app := ctx.Param("app")
	appLabel := config.LogAppLabel()
	options := config.ConvertRequestQueryParams(ctx.Request.URL.Query())
	options[appLabel] = app

	results, err := s.Service.Search(keyword, options, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
}

// Context search log context
func (s *Search) Context(ctx *gin.Context) {
	app := ctx.Param("app")
	appLabel := config.LogAppLabel()
	options := config.ConvertRequestQueryParams(ctx.Request.URL.Query())
	options[appLabel] = app
	results, err := s.Service.Context(options, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogContextError, err))
		return
	}

	utils.Ok(ctx, results)
	return
}

// Everything is sweet API for get log filter by gived conditions
// key is what you want return and must in config.logOptionalLabels
func (s *Search) Everything(ctx *gin.Context) {
	key := ctx.Param("key")
	options := config.ConvertRequestQueryParams(ctx.Request.URL.Query())
	results, err := s.Service.Everything(key, options, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	utils.Ok(ctx, results)
	return
}
