package service

import (
	"encoding/json"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

const (
	PROMETHEUS_INDEX = ".dataman-prometheus"
	PROMETHEUS_TYPE  = "dataman-prometheus"
)

func (s *SearchService) SavePrometheus(pro models.CommonLabels) error {
	pro.CreateTime = time.Now().Format(time.RFC3339Nano)
	_, err := s.ESClient.Index().
		Index(PROMETHEUS_INDEX).
		Type(PROMETHEUS_TYPE).
		BodyJson(pro).
		Do()
	return err
}

func (s *SearchService) GetPrometheus(page models.Page) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	result, err := s.ESClient.Search().
		Index(PROMETHEUS_INDEX).
		Type(PROMETHEUS_TYPE).
		From(page.PageFrom).
		Size(page.PageSize).
		Pretty(true).
		Do()

	if err != nil {
		return nil, err
	}

	var cls []models.CommonLabels
	for _, hit := range result.Hits.Hits {
		var pro models.Prometheus
		if json.Unmarshal(*hit.Source, &pro) == nil {
			cls = append(cls, pro.CLs)
		}
	}

	data["results"] = cls
	data["count"] = result.Hits.TotalHits

	return data, nil
}