package datastore

import (
	"errors"
	"strconv"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateLogAlertEvent(event *models.LogAlertEvent) error {
	return db.Save(event).Error
}

func (db *datastore) GetLogAlertEvents(options map[string]interface{}, page models.Page) (map[string]interface{}, error) {
	var (
		count  int
		events []*models.LogAlertEvent
	)

	query := db.Table("log_alert_events")

	if start, ok := options["start"]; ok {
		timeStamp, err := strconv.ParseInt(start.(string), 10, 64)
		if err != nil {
			return nil, err
		}

		startTime := time.Unix(timeStamp, 0)
		query = query.Where("log_time >= ?", startTime)

		delete(options, "start")
	}

	if end, ok := options["end"]; ok {
		timeStamp, err := strconv.ParseInt(end.(string), 10, 64)
		if err != nil {
			return nil, err
		}

		endTime := time.Unix(timeStamp, 0)
		query = query.Where("log_time <= ?", endTime)

		delete(options, "end")
	}

	if err := query.Where(options).Find(&events).Count(&count).Error; err != nil {
		return nil, err
	}

	if err := query.Where(options).Offset(page.PageFrom).Limit(page.PageSize).Order("log_time desc").
		Scan(&events).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{"count": count, "events": events}, nil
}

func (db *datastore) DeleteLogAlertEvents(start, end string) error {
	if start == "" || end == "" {
		return errors.New("interval start or end time")
	}

	query := db.Table("log_alert_events")

	startTimestamp, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		return err
	}

	startTime := time.Unix(startTimestamp, 0)
	query = query.Where("log_time >= ?", startTime)

	endTimestamp, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		return err
	}

	endTime := time.Unix(endTimestamp, 0)
	query = query.Where("log_time <= ?", endTime)

	return query.Delete(&models.LogAlertEvent{}).Error
}
