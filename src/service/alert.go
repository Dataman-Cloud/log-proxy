package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"gopkg.in/olivere/elastic.v3"
)

const (
	ALERT_INDEX = ".dataman-alerts"
	ALERT_TYPE  = "dataman-alerts"

	ALERT_HISTORY_INDEX = ".dataman-keyword-history"
	ALERT_HISTORY_TYPE  = "dataman-keyword-history"
)

func (s *SearchService) CreateAlert(alert *models.Alert) error {
	alert.CreateTime = time.Now()
	_, err := s.ESClient.Index().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		BodyJson(alert).
		Do()
	return err
}

func (s *SearchService) DeleteAlert(id string) error {
	_, err := s.ESClient.Delete().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		Id(id).
		Do()
	return err
}

func (s *SearchService) GetAlerts(page models.Page) ([]models.Alert, error) {
	var results []models.Alert

	result, err := s.ESClient.Search().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		From(page.PageFrom).
		Size(page.PageSize).
		Pretty(true).
		Do()
	if err != nil {
		return results, err
	}

	for _, hit := range result.Hits.Hits {
		data := models.Alert{
			Id: hit.Id,
		}

		json.Unmarshal(*hit.Source, &data)
		results = append(results, data)
	}

	return results, nil
}

func (s *SearchService) UpdateAlert(alert *models.Alert) error {
	_, err := s.ESClient.Update().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		Id(alert.Id).
		Doc(alert).
		Do()
	return err
}

func (s *SearchService) GetAlertCondition() []models.Alert {
	var alerts []models.Alert
	result, err := s.ESClient.Search().
		Index(".dataman-alerts").
		Type("dataman-alerts").
		Pretty(true).
		Do()
	if err != nil {
		return alerts
	}

	for _, hit := range result.Hits.Hits {
		alert := models.Alert{
			Id: hit.Id,
		}
		json.Unmarshal(*hit.Source, &alert)
		alerts = append(alerts, alert)
	}

	return alerts
}

func (s *SearchService) ExecuteAlert(alert models.Alert) {
	s.Maf.Lock()
	defer s.Maf.Unlock()

	t, ok := s.AlertFlag[alert.Id]
	if ok && !time.Now().After(t.Add(+time.Duration(alert.Period)*time.Minute)) {
		return
	}

	query := elastic.NewBoolQuery().
		Filter(elastic.NewRangeQuery("logtime.timestamp").Gte(fmt.Sprintf("now-%dm", alert.Period)).Lte("now")).
		Must(elastic.NewQueryStringQuery("message:" + alert.Keyword).AnalyzeWildcard(true))

	clusterName := strings.Split(alert.AppId, "-")[0]
	result, err := s.ESClient.Search().
		Index(fmt.Sprintf("dataman-%s-*", clusterName)).
		Type("dataman-"+alert.AppId).
		Query(query).
		Aggregation("tasks", elastic.NewTermsAggregation().Field("taskid").OrderByCountDesc()).
		Pretty(true).
		SearchType("count").
		Do()

	if err != nil {
		return
	}
	s.AlertFlag[alert.Id] = time.Now()

	s.CreateKeywordAlertInfo(models.KeywordAlertHistory{
		AppId:      alert.AppId,
		Keyword:    alert.Keyword,
		Count:      result.Hits.TotalHits,
		CreateTime: time.Now(),
	})
}

func (s *SearchService) CreateKeywordAlertInfo(info models.KeywordAlertHistory) {
	info.CreateTime = time.Now()
	s.ESClient.Index().
		Index(ALERT_HISTORY_INDEX).
		Type(ALERT_HISTORY_TYPE).
		BodyJson(info).
		Do()
}

func (s *SearchService) GetKeywordAlertHistory() (map[string]interface{}, error) {
	result, err := s.ESClient.Search().
		Index(ALERT_HISTORY_INDEX).
		Type(ALERT_HISTORY_TYPE).
		Pretty(true).
		Do()
	if err != nil {
		return nil, err
	}

	var results []models.KeywordAlertHistory
	for _, hit := range result.Hits.Hits {
		var kh models.KeywordAlertHistory
		json.Unmarshal(*hit.Source, &kh)
		results = append(results, kh)
	}
	return map[string]interface{}{
		"results": results,
		"count":   result.Hits.TotalHits,
	}, nil
}
