package service

import (
	"encoding/json"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

const (
	EVENT_INDEX = ".dataman-event"
	EVENT_TYPE  = "dataman-event"

	TASK_FAILED = "TASK_FAILED"
)

func (s *SearchService) SaveFaildEvent(data []byte) error {
	var event models.StatusUpdate
	err := json.Unmarshal(data, &event)
	if err != nil {
		return err
	}

	if event.TaskStatus != TASK_FAILED {
		return nil
	}
	event.CreateTime = time.Now().Format(time.RFC3339Nano)

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
