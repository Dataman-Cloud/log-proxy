package datastore

import (
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

	if err := db.Table("log_alert_events").Where(options).Find(&events).Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Table("log_alert_events").Where(options).Offset(page.PageFrom).Limit(page.PageSize).Order("log_time desc").
		Scan(&events).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{"count": count, "events": events}, nil
}
