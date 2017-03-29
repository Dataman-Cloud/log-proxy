package datastore

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateOrIncreaseEvent(event *models.Event) error {
	var result models.Event
	notfound := db.Where("ack = ? AND severity = ? AND cluster = ? AND app = ? AND task = ? AND container_id = ? AND container_name =? AND alert_name = ?",
		false,
		event.Severity,
		event.Cluster,
		event.App,
		event.Task,
		event.ContainerID,
		event.ContainerName,
		event.AlertName,
	).First(&result).RecordNotFound()
	if notfound {
		event.Count = 1
		return db.Create(event).Error
	}

	result.Count++
	result.Description = event.Description
	result.Summary = event.Summary
	result.Value = event.Value
	return db.Save(&result).Error

}

func (db *datastore) AckEvent(ID int, cluster, app string) error {
	var result models.Event
	if db.Where("id = ?", ID).Where("cluster = ? AND app = ?", cluster, app).First(&result).RecordNotFound() {
		return fmt.Errorf("Alert Event id=%d, cluster=%s or app=%s not found", ID, cluster, app)
	}
	result.Ack = true
	return db.Save(&result).Error
}

func (db *datastore) ListEvents(page models.Page, options map[string]interface{}) (map[string]interface{}, error) {
	var (
		result []*models.Event
		count  int
	)
	query := db.Table("events")

	if start, ok := options["start"]; ok {
		timeStamp, err := strconv.ParseInt(start.(string), 10, 64)
		if err != nil {
			return nil, err
		}

		startTime := time.Unix(timeStamp, 0)
		query = query.Where("updated_at >= ?", startTime)

		delete(options, "start")
	}

	if end, ok := options["end"]; ok {
		timeStamp, err := strconv.ParseInt(end.(string), 10, 64)
		if err != nil {
			return nil, err
		}

		endTime := time.Unix(timeStamp, 0)
		query = query.Where("updated_at <= ?", endTime)

		delete(options, "end")
	}

	err := query.
		Where(options).
		Find(&result).
		Count(&count).
		Error
	if err != nil {
		return nil, err
	}

	err = query.Where(options).
		Offset(page.PageFrom).
		Limit(page.PageSize).
		Order("updated_at desc").
		Find(&result).
		Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"count": count, "events": result}, nil
}
