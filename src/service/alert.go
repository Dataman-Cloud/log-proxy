package service

import (
	"encoding/json"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

const (
	ALERT_INDEX = ".dataman-alerts"
	ALERT_TYPE  = "dataman-alerts"
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

func (s *SearchService) GetAlerts() ([]models.Alert, error) {
	var results []models.Alert

	result, err := s.ESClient.Search().
		Index(ALERT_INDEX).
		Type(ALERT_TYPE).
		From(s.PageFrom).
		Size(s.PageSize).
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
