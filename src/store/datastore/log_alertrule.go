package datastore

import (
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateLogAlertRule(rule *models.LogAlertRule) error {
	var result models.LogAlertRule
	notfound := db.Where("log_alert_rules.app = ? AND log_alert_rules.source = ? AND log_alert_rules.keyword = ? "+
		"AND log_alert_rules.groupname = ?",
		rule.App, rule.Source, rule.Keyword, rule.Group).
		First(&result).
		RecordNotFound()
	if !notfound {
		return errors.New("The rule has been existed in Database")
	}

	return db.Save(rule).Error
}

func (db *datastore) UpdateLogAlertRule(rule *models.LogAlertRule) error {
	var result models.LogAlertRule
	notfound := db.Where("log_alert_rules.id = ?", rule.ID).
		First(&result).
		RecordNotFound()
	if notfound {
		return errors.New("The rule not found in Database")
	}

	return db.Model(rule).
		Where("log_alert_rules.id = ?", rule.ID).
		Omit("app, source").
		Updates(rule).Error
}

func (db *datastore) DeleteLogAlertRule(ID string) error {
	return db.Where("log_alert_rules.id = ?", ID).Delete(&models.LogAlertRule{}).Error
}

func (db *datastore) GetLogAlertRule(ID string) (models.LogAlertRule, error) {
	var result models.LogAlertRule
	err := db.Table("log_alert_rules").Where("ID = ?", ID).Scan(&result).Error
	return result, err
}

func (db *datastore) GetLogAlertRules(opts map[string]interface{}, page models.Page) (map[string]interface{}, error) {
	var (
		count int
		rules []*models.LogAlertRule
	)
	groups := opts["groups"]
	delete(opts, "groups")

	if err := db.Table("log_alert_rules").
		Where(opts).
		Where("groupname in (?)", groups).
		Find(&rules).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	if err := db.Table("log_alert_rules").
		Where(opts).
		Where("groupname in (?)", groups).
		Offset(page.PageFrom).
		Limit(page.PageSize).
		Scan(&rules).Error; err != nil {
		return nil, err

	}

	return map[string]interface{}{"count": count, "rules": rules}, nil
}

func (db *datastore) CreateLogAlertEvent(event *models.LogAlertEvent) error {
	return db.Save(event).Error
}

func (db *datastore) GetLogAlertEvents(opts map[string]interface{}, page models.Page) (map[string]interface{}, error) {
	var (
		count  int
		events []*models.LogAlertEvent
	)

	query := db.Table("log_alert_events")
	if page.RangeFrom != nil {
		query = query.Where("log_time >= ?", page.RangeFrom)
	}

	if page.RangeTo != nil {
		query = query.Where("log_time <= ?", page.RangeTo)
	}

	if err := query.Where(opts).Find(&events).Count(&count).Error; err != nil {
		return nil, err
	}

	if err := query.Where(opts).Offset(page.PageFrom).Limit(page.PageSize).Order("log_time desc").
		Scan(&events).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{"count": count, "events": events}, nil
}

func (db *datastore) GetLogAlertEvent(ID string) (*models.LogAlertEvent, error) {
	var event models.LogAlertEvent
	err := db.Table("log_alert_events").Where("id = ?", ID).Scan(&event).Error
	return &event, err
}

func (db *datastore) AckLogAlertEvent(ID string) error {
	return db.Table("log_alert_events").Where("id = ?", ID).Updates(map[string]interface{}{"ack": true}).Error
}

func (db *datastore) GetLogAlertApps(opts map[string]interface{}, page models.Page) ([]*models.LogAlertApps, error) {
	var apps []*models.LogAlertApps

	query := db.Table("log_alert_events").Select("DISTINCT(log_alert_events.app)")
	if page.RangeFrom != nil {
		query = query.Where("log_time >= ?", page.RangeFrom)
	}

	if page.RangeTo != nil {
		query = query.Where("log_time <= ?", page.RangeTo)
	}

	if err := query.Where(opts).Scan(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}
