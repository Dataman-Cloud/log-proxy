package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v3"
)

type SearchService struct {
	RangeFrom string
	RangeTo   string
	PageSize  int
	PageFrom  int
	ESClient  *elastic.Client
}

type SearchResult struct {
	Message []string `json:"message"`
}

func NewSearchService() *SearchService {
	var ofs []elastic.ClientOptionFunc
	ofs = append(ofs, elastic.SetURL(strings.Split(config.GetConfig().ES_URL, ",")...))
	if config.GetConfig().SEARCH_DEBUG {
		ofs = append(ofs, elastic.SetTraceLog(log.StandardLogger()))
	}
	client, err := elastic.NewClient(ofs...)
	if err != nil {
		log.Errorf("new elastic client error: %v", err)
		return nil
	}

	return &SearchService{
		RangeFrom: "now-15m",
		RangeTo:   "now",
		PageFrom:  0,
		PageSize:  100,
		ESClient:  client,
	}
}

func (s *SearchService) Applications() (map[string]int64, error) {
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(s.RangeFrom).Lte(s.RangeTo).Format("epoch_millis"))

	apps := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation("apps", elastic.
			NewTermsAggregation().
			Field("appid").
			OrderByCountDesc()).
		Pretty(true).
		Do()
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

func (s *SearchService) Tasks(appName string) (map[string]int64, error) {
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(s.RangeFrom).Lte(s.RangeTo).Format("epoch_millis")).
		Must(elastic.NewTermQuery("appid", appName))

	tasks := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		Type("dataman-"+appName).
		Query(bquery).
		SearchType("count").
		Aggregation("tasks", elastic.
			NewTermsAggregation().
			Field("taskid").
			OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil {
		log.Errorf("get app %s tasks error: %v", appName, err)
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

func (s *SearchService) Paths(appName, taskId string) (map[string]int64, error) {
	paths := make(map[string]int64)

	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(s.RangeFrom).Lte(s.RangeTo).Format("epoch_millis")).
		Must(elastic.NewTermQuery("appid", appName), elastic.NewTermQuery("taskid", taskId))

	result, err := s.ESClient.Search().
		Index("dataman-*").
		Type("dataman-"+appName).
		Query(bquery).
		SearchType("count").
		Aggregation("paths", elastic.
			NewTermsAggregation().
			Field("path").
			OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil {
		log.Errorf("get app %s paths error: %v", appName, err)
		return paths, err
	}

	agg, found := result.Aggregations.Terms("paths")
	if !found {
		return paths, nil
	}

	for _, bucket := range agg.Buckets {
		paths[fmt.Sprint(bucket.Key)] = bucket.DocCount
	}

	return paths, nil
}

func (s *SearchService) Search(appid, taskid, source, keyword string) (map[string]interface{}, error) {
	var data map[string]interface{}

	var querys []elastic.Query
	if taskid != "" {
		querys = append(querys, elastic.NewTermsQuery("taskid", strings.Split(taskid, ",")))
	}
	if source != "" {
		querys = append(querys, elastic.NewTermsQuery("path", strings.Split(source, ",")))
	}
	if keyword != "" {
		querys = append(querys, elastic.NewQueryStringQuery("message:"+keyword).AnalyzeWildcard(true))
	}
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(s.RangeFrom).Lte(s.RangeTo).Format("epoch_millis")).
		Must(querys...)

	result, err := s.ESClient.Search().
		Index("dataman-"+strings.Split(appid, "-")[0]+"-"+utils.ParseDate(s.RangeFrom, s.RangeTo)).
		Type("dataman-"+appid).
		Query(bquery).
		Highlight(elastic.NewHighlight().Field("message").PreTags(`@dataman-highlighted-field@`).PostTags(`@/dataman-highlighted-field@`)).
		Sort("offset", true).From(s.PageFrom).Size(s.PageSize).Pretty(true).IgnoreUnavailable(true).Do()

	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, hit := range result.Hits.Hits {
		data := make(map[string]interface{})
		json.Unmarshal(*hit.Source, &data)
		if len(hit.Highlight["message"]) > 0 {
			data["message"] = hit.Highlight["message"][0]
		}
		results = append(results, data)
	}
	data["results"] = results
	data["count"] = result.Hits.TotalHits

	return data, nil
}
