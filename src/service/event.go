package service

import (
	"encoding/json"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

const (
	EVENT_INDEX = ".dataman-event"
	EVENT_TYPE  = "dataman-event"
)

func (s *SearchService) SaveFaildEvent(data []byte) error {
	var event models.StatusUpdate
	err := json.Unmarshal(data, &event)
	if err != nil {
		return err
	}

	_, err = s.ESClient.Index().
		Index(EVENT_INDEX).
		Type(EVENT_TYPE).
		BodyJson(event).
		Do()
	return err
}

func (s *SearchService) GetEvents(page models.Page) (map[string]interface{}, error) {
	var data map[string]interface{}

	result, err := s.ESClient.Search().
		Index(EVENT_INDEX).
		Type(EVENT_TYPE).
		From(page.PageFrom).
		Size(page.PageSize).
		Pretty(true).
		Do()

	if err != nil {
		return nil, err
	}

	var events []models.StatusUpdate
	for _, hit := range result.Hits.Hits {
		var event models.StatusUpdate
		if json.Unmarshal(*hit.Source, &event) == nil {
			events = append(events, event)
		}
	}

	data["results"] = events
	data["count"] = result.Hits.TotalHits

	return data, nil
}
