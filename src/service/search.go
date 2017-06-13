package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils/esclient"

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
func NewEsService(urls []string) *SearchService {
	client, err := esclient.New(urls)
	if err != nil {
		log.Errorf("new elastic client error: %v", err)
	}
	return &SearchService{
		ESClient:  client,
		AlertFlag: make(map[string]time.Time),
		Maf:       new(sync.Mutex),
	}
}

func (s *SearchService) resetESClient() error {
	var err error
	if s.ESClient == nil {
		s.ESClient, err = esclient.GetClient(strings.Split(config.GetConfig().EsURL, ","))
		if err != nil {
			return err
		}
	}
	return nil
}

// Applications get all applications
func (s *SearchService) Applications(page models.Page) (map[string]int64, error) {
	var err error
	err = s.resetESClient()
	if err != nil {
		return nil, err
	}
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

	if err != nil {
		switch errType := err.(type) {
		case *url.Error:
			err = s.resetESClient()
			if err != nil {
				return nil, err
			}
		case *elastic.Error:
			if errType.Status == http.StatusNotFound {
				return nil, nil
			}
		}
	}

	if err != nil {
		log.Errorf("get applications error: %v", err)
		return apps, err
	}

	if result == nil {
		return nil, fmt.Errorf("Get the null from elasicsearch")
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
	var err error
	err = s.resetESClient()
	if err != nil {
		return nil, err
	}

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

	if result == nil {
		return nil, fmt.Errorf("Get the null from elasicsearch")
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
	var err error
	err = s.resetESClient()
	if err != nil {
		return nil, err
	}

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
	var err error
	err = s.resetESClient()
	if err != nil {
		return nil, err
	}

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
func (s *SearchService) Search(keyword string, opts map[string]interface{}, page models.Page) (map[string]interface{}, error) {
	var err error
	err = s.resetESClient()
	if err != nil {
		return nil, err
	}

	sort := false
	var querys []elastic.Query

	messageLabel := config.LogMessageLabel()
	if keyword != "" {
		querys = append(querys, elastic.NewQueryStringQuery(messageLabel+":"+keyword).AnalyzeWildcard(true))
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
	var err error
	err = s.resetESClient()
	if err != nil {
		return nil, err
	}

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
	var err error
	err = s.resetESClient()
	if err != nil {
		return nil, err
	}

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
