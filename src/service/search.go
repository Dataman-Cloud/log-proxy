package service

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

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

func (s *SearchService) Clusters(page models.Page) (map[string]int64, error) {
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis"))

	clusters := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation("clusters", elastic.NewTermsAggregation().Field("DM_VCLUSTER").Size(0).OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		log.Errorf("query clusters got error: %v", err)
		return clusters, err
	}

	agg, found := result.Aggregations.Terms("clusters")
	if !found {
		return clusters, nil
	}

	for _, bucket := range agg.Buckets {
		clusters[fmt.Sprint(bucket.Key)] = bucket.DocCount
	}
	return clusters, nil
}

// Applications get all applications
func (s *SearchService) Applications(cluster string, page models.Page) (map[string]int64, error) {
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(elastic.NewTermQuery("DM_VCLUSTER", cluster))

	apps := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation("apps", elastic.NewTermsAggregation().Field("DM_APP_ID").Size(0).OrderByCountDesc()).
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

func (s *SearchService) Slots(cluster, app string, page models.Page) (map[string]int64, error) {
	bQuery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(elastic.NewTermsQuery("DM_VCLUSTER", cluster), elastic.NewTermQuery("DM_APP_ID", app))

	slots := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bQuery).
		Aggregation("slots", elastic.NewTermsAggregation().Field("DM_SLOT_INDEX").Size(0).OrderByCountDesc()).
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
func (s *SearchService) Tasks(cluster, app, slot string, page models.Page) (map[string]int64, error) {
	var querys []elastic.Query
	querys = append(querys, elastic.NewTermQuery("DM_VCLUSTER", cluster))
	querys = append(querys, elastic.NewTermQuery("DM_APP_ID", app))
	querys = append(querys, elastic.NewTermQuery("DM_SLOT_INDEX", slot))
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(querys...)

	//Index("dataman-*").
	tasks := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation("tasks", elastic.NewTermsAggregation().Field("DM_TASK_ID").Size(0).OrderByCountDesc()).
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

// Paths get application paths
func (s *SearchService) Paths(cluster, user, appName, taskID string, page models.Page) (map[string]int64, error) {
	paths := make(map[string]int64)
	var querys []elastic.Query
	querys = append(querys, elastic.NewTermQuery("cluster", cluster))
	querys = append(querys, elastic.NewTermQuery("user", user))
	querys = append(querys, elastic.NewTermQuery("app", appName))
	if taskID != "" {
		querys = append(querys, elastic.NewTermsQuery("task", utils.ParseTask(taskID)))
	}

	bquery := elastic.NewBoolQuery().
		Filter(
			elastic.NewRangeQuery("logtime.timestamp").
				Gte(page.RangeFrom).
				Lte(page.RangeTo).Format("epoch_millis")).
		Must(querys...)

	result, err := s.ESClient.Search().
		Index("dataman-*-"+utils.ParseDate(page.RangeFrom, page.RangeTo)).
		Type("dataman-"+user+"-"+appName).
		Query(bquery).
		SearchType("count").
		Aggregation("paths", elastic.
			NewTermsAggregation().
			Field("path").
			Size(0).
			OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

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

// Search search log by condition
func (s *SearchService) Search(cluster, user, app, task, source, keyword string, page models.Page) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	sort := false
	var querys []elastic.Query
	querys = append(querys, elastic.NewTermQuery("cluster", cluster))
	querys = append(querys, elastic.NewTermQuery("user", user))
	if task != "" {
		querys = append(querys, elastic.NewTermsQuery("task", utils.ParseTask(task)))
	}
	if source != "" {
		querys = append(querys, elastic.NewTermsQuery("path", strings.Split(source, ",")))
	}
	if keyword != "" {
		sort = true
		querys = append(querys, elastic.NewQueryStringQuery("message:"+keyword).AnalyzeWildcard(true))
	}

	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(querys...)

	result, err := s.ESClient.Search().
		Index("dataman-"+cluster+"-"+utils.ParseDate(page.RangeFrom, page.RangeTo)).
		Type("dataman-"+user+"-"+app).
		Query(bquery).
		Sort("logtime.sort", sort).
		Highlight(elastic.NewHighlight().Field("message").PreTags(`@dataman-highlighted-field@`).PostTags(`@/dataman-highlighted-field@`)).
		From(page.PageFrom).Size(page.PageSize).Pretty(true).IgnoreUnavailable(true).Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, hit := range result.Hits.Hits {
		data := make(map[string]interface{})
		json.Unmarshal(*hit.Source, &data)
		if len(hit.Highlight["message"]) > 0 {
			str := html.EscapeString(hit.Highlight["message"][0])
			str = strings.Replace(str, "@dataman-highlighted-field@", "<mark>", -1)
			str = strings.Replace(str, "@/dataman-highlighted-field@", "</mark>", -1)
			data["message"] = str
		}
		results = append(results, data)
	}
	data["results"] = results
	data["count"] = result.Hits.TotalHits

	return data, nil
}

// Context search log context
func (s *SearchService) Context(cluster, user, app, task, source, timestamp string, page models.Page) ([]map[string]interface{}, error) {
	offset, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	var results []map[string]interface{}

	if page.PageFrom == 0 {
		bquery := elastic.NewBoolQuery().
			Filter(elastic.NewRangeQuery("offset").Lt(offset)).
			Must(
				elastic.NewTermQuery("cluster", cluster),
				elastic.NewTermQuery("user", user),
				elastic.NewTermQuery("app", app),
				elastic.NewTermQuery("task", task),
				elastic.NewTermQuery("path", source),
			)

		result, err := s.ESClient.Search().
			Index("dataman-"+cluster+"-"+time.Unix(offset/1e9, 0).Format("2006-01-02")).
			Type("dataman-"+user+"-"+app).
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
		Filter(elastic.NewRangeQuery("offset").Gte(offset)).
		Must(elastic.NewTermQuery("app", app), elastic.NewTermQuery("task", task), elastic.NewTermQuery("path", source))

	result, err := s.ESClient.Search().
		Index("dataman-"+cluster+"-"+time.Unix(offset/1e9, 0).Format("2006-01-02")).
		Type("dataman-"+user+"-"+app).
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
