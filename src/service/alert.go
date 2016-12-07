package service

import (
	"encoding/json"
	"net/http"
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
	alert.CreateTime = time.Now().Format(time.RFC3339Nano)
	_, err := s.ESClient.Index().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		BodyJson(alert).
		Do()
	s.ESClient.Flush()
	return err
}

func (s *SearchService) DeleteAlert(id string) error {
	_, err := s.ESClient.Delete().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		Id(id).
		Do()
	s.ESClient.Flush()
	return err
}

func (s *SearchService) GetAlerts(page models.Page) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	var results []models.Alert

	result, err := s.ESClient.Search().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		From(page.PageFrom).
		Size(page.PageSize).
		Sort("createtime.timestamp", true).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	for _, hit := range result.Hits.Hits {
		data := models.Alert{
			Id: hit.Id,
		}

		json.Unmarshal(*hit.Source, &data)
		results = append(results, data)
	}

	m["results"] = results
	m["count"] = result.Hits.TotalHits

	return m, nil
}

func (s *SearchService) GetAlert(id string) (models.Alert, error) {
	var alert models.Alert

	result, err := s.ESClient.Get().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		Id(id).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return alert, nil
	}

	if err != nil {
		return alert, err
	}

	err = json.Unmarshal(*result.Source, &alert)
	if err != nil {
		return alert, err
	}
	alert.Id = result.Id

	return alert, nil
}

func (s *SearchService) UpdateAlert(alert *models.Alert) error {
	_, err := s.ESClient.Update().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		Id(alert.Id).
		Doc(alert).
		Do()
	s.ESClient.Flush()
	return err
}

func (s *SearchService) GetAlertCondition() []models.Alert {
	var alerts []models.Alert
	result, err := s.ESClient.Search().
		Index(".dataman-alerts").
		Type("dataman-alerts").
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil
	}

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
