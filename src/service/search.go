package service

import (
	"fmt"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/config"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v3"
)

type SearchService struct {
	ESClient *elastic.Client
}

func GetSearchService() *SearchService {
	client, err := elastic.NewClient(elastic.SetURL(strings.Split(config.GetConfig().ES_URL, ",")...))
	if err != nil {
		log.Errorf("new elastic client error: %v", err)
		return nil
	}

	return &SearchService{
		ESClient: client,
	}
}

func (s *SearchService) Applications() (map[string]int64, error) {
	apps := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		Query(elastic.NewMatchAllQuery()).
		SearchType("count").
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
	tasks := make(map[string]int64)
	result, err := s.ESClient.Search().
		Index("dataman-*").
		Type("dataman-"+appName).
		Query(elastic.NewTermQuery("appid", appName)).
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
	result, err := s.ESClient.Search().
		Index("dataman-*").
		Type("dataman-"+appName).
		Query(elastic.NewTermQuery("appid", appName)).
		Query(elastic.NewTermQuery("taskid", taskId)).
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

func (s *SearchService) Search() {
}
