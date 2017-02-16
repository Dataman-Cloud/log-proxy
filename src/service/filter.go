package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"gopkg.in/olivere/elastic.v3"
)

const (
	// KWFilterIndex keyword filter es index
	KWFilterIndex = ".dataman-keyword-filter"
	// KWFilterType keyword filter es type
	KWFilterType = "dataman-keyword-filter"
)

// CreateFilter create filter
func (s *SearchService) CreateFilter(filter *models.KWFilter) error {
	filter.CreateTime = time.Now().Format(time.RFC3339Nano)
	_, err := s.ESClient.Index().
		Index(KWFilterIndex).
		Type(KWFilterType).
		BodyJson(filter).
		Do()
	s.ESClient.Flush()
	return err
}

// DeleteFilter delete filter by id
func (s *SearchService) DeleteFilter(id string) error {
	_, err := s.ESClient.Delete().
		Index(KWFilterIndex).
		Type(KWFilterType).
		Id(id).
		Do()
	s.ESClient.Flush()
	return err
}

// GetFilters get all filters
func (s *SearchService) GetFilters(page models.Page) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	var results []models.KWFilter

	result, err := s.ESClient.Search().
		Index(KWFilterIndex).
		Type(KWFilterType).
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
		data := models.KWFilter{
			ID: hit.Id,
		}

		json.Unmarshal(*hit.Source, &data)
		results = append(results, data)
	}

	m["results"] = results
	m["count"] = result.Hits.TotalHits

	return m, nil
}

// GetFilter get filter by id
func (s *SearchService) GetFilter(id string) (models.KWFilter, error) {
	var filter models.KWFilter

	result, err := s.ESClient.Get().
		Index(KWFilterIndex).
		Type(KWFilterType).
		Id(id).
		Do()

	if err != nil && err.(*elastic.Error).Status == http.StatusNotFound {
		return filter, nil
	}

	if err != nil {
		return filter, err
	}

	err = json.Unmarshal(*result.Source, &filter)
	if err != nil {
		return filter, err
	}
	filter.ID = result.Id

	return filter, nil
}

// UpdateFilter update filter
func (s *SearchService) UpdateFilter(filter *models.KWFilter) error {
	_, err := s.ESClient.Update().
		Index(KWFilterIndex).
		Type(KWFilterType).
		Id(filter.ID).
		Doc(filter).
		Do()
	s.ESClient.Flush()
	return err
}
