package datastore

import (
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateLogAlertRule(rule *models.LogAlertRule) error {
	var result models.LogAlertRule
	notfound := db.Where("log_alert_rules.app = ? AND log_alert_rules.source = ? AND log_alert_rules.keyword = ? "+
		"AND log_alert_rules.user = ? AND log_alert_rules.group = ?",
		rule.App, rule.Source, rule.Keyword, rule.User, rule.Group).
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

func (db *datastore) GetLogAlertRules(group string, page models.Page) (map[string]interface{}, error) {
	var (
		count int
		rules []*models.LogAlertRule
	)

	if err := db.Table("log_alert_rules").Where("log_alert_rules = ?", group).Find(&rules).Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Table("log_alert_rules").Offset(page.PageFrom).Limit(page.PageSize).
		Scan(&rules).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{"count": count, "rules": rules}, nil
}
