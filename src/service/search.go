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

// Applications get all applications
func (s *SearchService) Applications(page models.Page) (map[string]int64, error) {
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis"))

	apps := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		SearchType("count").
		Query(bquery).
		Aggregation("apps", elastic.
			NewTermsAggregation().
			Field("appid").
			Size(0).
			OrderByCountDesc()).
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

// Tasks get application tasks
func (s *SearchService) Tasks(appName string, page models.Page) (map[string]int64, error) {
	bquery := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(page.RangeFrom).Lte(page.RangeTo).Format("epoch_millis")).
		Must(elastic.NewTermQuery("appid", appName))

	//Index("dataman-*").
	tasks := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*-"+utils.ParseDate(page.RangeFrom, page.RangeTo)).
		Type("dataman-"+appName).
		Query(bquery).
		SearchType("count").
		Aggregation("tasks", elastic.
			NewTermsAggregation().
			Field("taskid").
			Size(0).
			OrderByCountDesc()).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

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

// Paths get application paths
func (s *SearchService) Paths(clusterid, userid, appName, taskID string, page models.Page) (map[string]int64, error) {
	paths := make(map[string]int64)
	var querys []elastic.Query
	querys = append(querys, elastic.NewTermQuery("clusterid", clusterid))
	querys = append(querys, elastic.NewTermQuery("userid", userid))
	querys = append(querys, elastic.NewTermQuery("appid", appName))
	if taskID != "" {
		querys = append(querys, elastic.NewTermsQuery("taskid", utils.ParseTask(taskID)))
	}

	bquery := elastic.NewBoolQuery().
		Filter(
			elastic.NewRangeQuery("logtime.timestamp").
				Gte(page.RangeFrom).
				Lte(page.RangeTo).Format("epoch_millis")).
		Must(querys...)

	result, err := s.ESClient.Search().
		Index("dataman-*-"+utils.ParseDate(page.RangeFrom, page.RangeTo)).
		Type("dataman-"+appName).
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
func (s *SearchService) Search(clusterid, userid, appid, taskid, source, keyword string, page models.Page) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	sort := false
	var querys []elastic.Query
	querys = append(querys, elastic.NewTermQuery("clusterid", clusterid))
	querys = append(querys, elastic.NewTermQuery("userid", userid))
	if taskid != "" {
		querys = append(querys, elastic.NewTermsQuery("taskid", utils.ParseTask(taskid)))
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
		Index("dataman-"+clusterid+"-"+utils.ParseDate(page.RangeFrom, page.RangeTo)).
		Type("dataman-"+appid).
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
func (s *SearchService) Context(clusterid, userid, appid, taskid, source, timestamp string, page models.Page) ([]map[string]interface{}, error) {
	offset, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	var results []map[string]interface{}

	if page.PageFrom == 0 {
		bquery := elastic.NewBoolQuery().
			Filter(elastic.NewRangeQuery("offset").Lt(offset)).
			Must(
				elastic.NewTermQuery("clusterid", clusterid),
				elastic.NewTermQuery("userid", userid),
				elastic.NewTermQuery("appid", appid),
				elastic.NewTermQuery("taskid", taskid),
				elastic.NewTermQuery("path", source),
			)

		result, err := s.ESClient.Search().
			Index("dataman-"+clusterid+"-"+time.Unix(offset/1e9, 0).Format("2006-01-02")).
			Type("dataman-"+appid).
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
		Must(elastic.NewTermQuery("appid", appid), elastic.NewTermQuery("taskid", taskid), elastic.NewTermQuery("path", source))

	result, err := s.ESClient.Search().
		Index("dataman-"+strings.Split(appid, "-")[0]+"-"+time.Unix(offset/1e9, 0).Format("2006-01-02")).
		Type("dataman-"+appid).
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
