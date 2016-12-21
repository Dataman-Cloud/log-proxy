package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"gopkg.in/olivere/elastic.v3"
)

const (
	// AlertIndex alert es index
	AlertIndex = ".dataman-alerts"
	// AlertType alert es type
	AlertType = "dataman-alerts"

	// AlertHistoryIndex alert history es index
	AlertHistoryIndex = ".dataman-keyword-history"
	// AlertHistoryType alert type es type
	AlertHistoryType = "dataman-keyword-history"
)

// CreateAlert create alert
func (s *SearchService) CreateAlert(alert *models.Alert) error {
	alert.CreateTime = time.Now().Format(time.RFC3339Nano)
	_, err := s.ESClient.Index().
		Index(AlertIndex).
		Type(AlertType).
		BodyJson(alert).
		Do()
	s.ESClient.Flush()
	return err
}

// DeleteAlert delete alert by id
func (s *SearchService) DeleteAlert(id string) error {
	_, err := s.ESClient.Delete().
		Index(AlertIndex).
		Type(AlertType).
		Id(id).
		Do()
	s.ESClient.Flush()
	return err
}

// GetAlerts get all alerts
func (s *SearchService) GetAlerts(page models.Page) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	var results []models.Alert

	result, err := s.ESClient.Search().
		Index(AlertIndex).
		Type(AlertType).
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
			ID: hit.Id,
		}

		json.Unmarshal(*hit.Source, &data)
		results = append(results, data)
	}

	m["results"] = results
	m["count"] = result.Hits.TotalHits

	return m, nil
}

// GetAlert get alert by id
func (s *SearchService) GetAlert(id string) (models.Alert, error) {
	var alert models.Alert

	result, err := s.ESClient.Get().
		Index(AlertIndex).
		Type(AlertType).
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
	alert.ID = result.Id

	return alert, nil
}

// UpdateAlert update alert
func (s *SearchService) UpdateAlert(alert *models.Alert) error {
	_, err := s.ESClient.Update().
		Index(AlertIndex).
		Type(AlertType).
		Id(alert.ID).
		Doc(alert).
		Do()
	s.ESClient.Flush()
	return err
}
