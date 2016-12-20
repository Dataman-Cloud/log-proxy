package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"gopkg.in/olivere/elastic.v3"
)

const (
	// PrometheusIndex prometheus index
	PrometheusIndex = ".dataman-prometheus"
	// PrometheusType prometheus type
	PrometheusType = "dataman-prometheus"
)

// SavePrometheus recive prometheus event and save to es
func (s *SearchService) SavePrometheus(pro map[string]interface{}) error {
	//pro.CreateTime = time.Now().Format(time.RFC3339Nano)
	pro["createtime"] = time.Now().Format(time.RFC3339Nano)
	_, err := s.ESClient.Index().
		Index(PrometheusIndex).
		Type(PrometheusType).
		BodyJson(pro).
		Do()
	return err
}

// GetPrometheus get all prometheus
func (s *SearchService) GetPrometheus(page models.Page) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	result, err := s.ESClient.Search().
		Index(PrometheusIndex).
		Type(PrometheusType).
		From(page.PageFrom).
		Size(page.PageSize).
		Sort("createtime.timestamp", false).
		Pretty(true).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var cls []map[string]interface{}
	for _, hit := range result.Hits.Hits {
		var pro map[string]interface{}
		if json.Unmarshal(*hit.Source, &pro) == nil {
			pro["id"] = hit.Id
			cls = append(cls, pro)
		}
	}

	data["results"] = cls
	data["count"] = result.Hits.TotalHits

	return data, nil
}

// GetPrometheu get prometheus by id
func (s *SearchService) GetPrometheu(id string) (map[string]interface{}, error) {
	result, err := s.ESClient.Get().
		Index(PrometheusIndex).
		Type(PrometheusType).
		Id(id).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(*result.Source, &m)

	return m, err
}
