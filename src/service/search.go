package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v3"
)

// SearchService search service client
type SearchService struct {
	ESClient  *elastic.Client
	AlertFlag map[string]time.Time
	Maf       *sync.Mutex
}

// SearchResult search log result
type SearchResult struct {
	Message []string `json:"message"`
}

// NewEsService new es search client
func NewEsService(url []string) *SearchService {
	var ofs []elastic.ClientOptionFunc
	ofs = append(ofs, elastic.SetURL(url...))
	ofs = append(ofs, elastic.SetMaxRetries(3))
	if config.GetConfig().SearchDebug {
		ofs = append(ofs, elastic.SetTraceLog(log.StandardLogger()))
	}
	client, err := elastic.NewClient(ofs...)
	if err != nil {
		log.Errorf("new elastic client error: %v", err)
		return nil
	}

	return &SearchService{
		ESClient:  client,
		AlertFlag: make(map[string]time.Time),
		Maf:       new(sync.Mutex),
	}
}

// Applications get all applications
func (s *SearchService) Applications(page models.Page) (map[string]int64, error) {
	appLabel := config.LogAppLabel()
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis"))

	apps := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation("apps", elastic.NewTermsAggregation().Field(appLabel).Size(0).OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		log.Errorf("get applications error: %v", err)
		return apps, err
	}

	agg, found := result.Aggregations.Terms("apps")
	if !found {
		return apps, nil
	}

	for _, bucket := range agg.Buckets {
		apps[fmt.Sprint(bucket.Key)] = bucket.DocCount
	}
	return apps, nil
}

func (s *SearchService) Slots(app string, page models.Page) (map[string]int64, error) {
	appLabel := config.LogAppLabel()
	slotLabel := config.LogSlotLabel()
	bQuery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(elastic.NewTermQuery(appLabel, app))

	slots := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bQuery).
		Aggregation("slots", elastic.NewTermsAggregation().Field(slotLabel).Size(0).OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil {
		if err.(*elastic.Error).Status == http.StatusNotFound {
			return slots, nil
		}
		return slots, err
	}

	agg, found := result.Aggregations.Terms("slots")
	if !found {
		return slots, nil
	}

	for _, bucket := range agg.Buckets {
		slots[fmt.Sprint(bucket.Key)] = bucket.DocCount
	}

	return slots, nil
}

// Tasks get application tasks
func (s *SearchService) Tasks(opts map[string]interface{}, page models.Page) (map[string]int64, error) {
	var querys []elastic.Query
	for k, v := range opts {
		querys = append(querys, elastic.NewTermQuery(k, v))
	}
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(querys...)

	//Index("dataman-*").
	tasks := make(map[string]int64)
	taskLabel := config.LogTaskLabel()
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation("tasks", elastic.NewTermsAggregation().Field(taskLabel).Size(0).OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return tasks, err
	}

	agg, found := result.Aggregations.Terms("tasks")
	if !found {
		return tasks, nil
	}

	for _, bucket := range agg.Buckets {
		tasks[fmt.Sprint(bucket.Key)] = bucket.DocCount
	}
	return tasks, nil
}

func (s *SearchService) Sources(opts map[string]interface{}, page models.Page) (map[string]int64, error) {
	sources := make(map[string]int64)
	var querys []elastic.Query
	for k, v := range opts {
		querys = append(querys, elastic.NewTermQuery(k, v))
	}

	sourceLabel := config.LogSourceLabel()
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(querys...)

	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation("sources", elastic.NewTermsAggregation().Field(sourceLabel).Size(0).OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return sources, err
	}

	agg, found := result.Aggregations.Terms("sources")
	if !found {
		return sources, nil
	}

	for _, bucket := range agg.Buckets {
		sources[fmt.Sprint(bucket.Key)] = bucket.DocCount
	}

	return sources, nil
}

// Search search log by condition
func (s *SearchService) Search(opts map[string]interface{}, page models.Page) (map[string]interface{}, error) {
	sort := false
	var querys []elastic.Query

	messageLabel := config.LogMessageLabel()
	keywordLabel := config.LogKeywordLabel()
	keyword, ok := opts[keywordLabel]
	if ok {
		sort = true
		keywordStr := keyword.(string)
		conjLabel := config.LogConjLabel()
		conj, ok := opts[conjLabel]
		if ok && strings.ToLower(conj.(string)) == "or" {
			keyword = strings.Join(strings.Split(keywordStr, " "), " OR ")
		} else {
			keyword = strings.Join(strings.Split(keywordStr, " "), " AND ")
		}

		if ok {
			delete(opts, conjLabel)
		}

		querys = append(querys, elastic.NewQueryStringQuery(messageLabel+keywordStr).AnalyzeWildcard(true))
		delete(opts, keywordLabel)
	}

	for k, v := range opts {
		querys = append(querys, elastic.NewTermQuery(k, v))
	}

	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(querys...)

	result, err := s.ESClient.Search().
		Index("dataman-*").
		Query(bquery).
		Sort("logtime.sort", sort).
		Highlight(elastic.NewHighlight().Field(messageLabel).PreTags(`@dataman-highlighted-field@`).PostTags(`@/dataman-highlighted-field@`)).
		From(page.PageFrom).
		Size(page.PageSize).
		Pretty(true).
		IgnoreUnavailable(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var logs []map[string]interface{}

	if result.Hits != nil {
		for _, hit := range result.Hits.Hits {
			logContent := make(map[string]interface{})
			json.Unmarshal(*hit.Source, &logContent)
			if len(hit.Highlight[messageLabel]) > 0 {
				str := html.EscapeString(hit.Highlight[messageLabel][0])
				str = strings.Replace(str, "@dataman-highlighted-field@", "<mark>", -1)
				str = strings.Replace(str, "@/dataman-highlighted-field@", "</mark>", -1)
				logContent[messageLabel] = str
			}

			logs = append(logs, logContent)
		}
	}

	data := make(map[string]interface{})
	data["results"] = logs
	data["count"] = result.Hits.TotalHits

	return data, nil
}

// Context search log context
func (s *SearchService) Context(opts map[string]interface{}, page models.Page) ([]map[string]interface{}, error) {
	offsetLabel := config.LogOffsetLabel()
	offset_, ok := opts[offsetLabel]
	if !ok {
		return nil, errors.New("offset should not be nil")
	}

	delete(opts, offsetLabel)
	offset, err := strconv.ParseInt(offset_.(string), 10, 64)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	var querys []elastic.Query
	for k, v := range opts {
		querys = append(querys, elastic.NewTermQuery(k, v))
	}

	if page.PageFrom == 0 {
		bquery := elastic.NewBoolQuery().
			Filter(elastic.NewRangeQuery(offsetLabel).Lt(offset)).
			Must(querys...)

		result, err := s.ESClient.Search().
			Index("dataman-*").
			Query(bquery).
			Sort("logtime.sort", false).
			From(0).
			Size(page.PageSize).
			Pretty(true).
			Do()

		if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}

		for i := len(result.Hits.Hits) - 1; i >= 0; i-- {
			data := make(map[string]interface{})
			json.Unmarshal(*result.Hits.Hits[i].Source, &data)
			results = append(results, data)
		}
	}

	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery(offsetLabel).Gte(offset)).
		Must(querys...)

	result, err := s.ESClient.Search().
		Index("dataman-*").
		Query(bquery).
		Sort("logtime.sort", true).
		From(page.PageFrom).
		Size(page.PageSize).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	for _, hit := range result.Hits.Hits {
		data := make(map[string]interface{})
		json.Unmarshal(*hit.Source, &data)
		results = append(results, data)
	}

	return results, nil
}

func (s *SearchService) Everything(key string, opts map[string]interface{}, page models.Page) (map[string]int64, error) {
	var querys []elastic.Query
	for k, v := range opts {
		querys = append(querys, elastic.NewTermQuery(k, v))
	}

	keyLabel, err := config.GetLogLabel(key)
	if err != nil {
		return nil, err
	}

	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(querys...)

	tasks := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation(key, elastic.NewTermsAggregation().Field(keyLabel).Size(0).OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return tasks, err
	}

	agg, found := result.Aggregations.Terms(key)
	if !found {
		return tasks, nil
	}

	for _, bucket := range agg.Buckets {
		tasks[fmt.Sprint(bucket.Key)] = bucket.DocCount
	}
	return tasks, nil
}
