package datastore

import (
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) ListAlertRules(page models.Page, name string) (map[string]interface{}, error) {
	var (
		rules []*models.Rule
		count int
		err   error
	)
	if name == "" {
		err = db.Table("rules").Find(&rules).Count(&count).Error
		err = db.Table("rules").Offset(page.PageFrom).Limit(page.PageSize).Find(&rules).Error
	} else {
		err = db.Table("rules").Where("name = ?", name).Scan(&rules).Count(&count).Error
		err = db.Table("rules").Where("name = ?", name).Offset(page.PageFrom).Limit(page.PageSize).Scan(&rules).Error
	}

	return map[string]interface{}{"count": count, "rules": rules}, err
}

func (db *datastore) GetAlertRules() ([]*models.Rule, error) {
	var (
		rules []*models.Rule
		err   error
	)
	err = db.Table("rules").Find(&rules).Error

	return rules, err
}

func (db *datastore) GetAlertRule(id uint64) (models.Rule, error) {
	var result models.Rule
	err := db.Table("rules").
		Where("id = ?", id).
		Scan(&result).Error
	return result, err
}

func (db *datastore) GetAlertRuleByName(name, alert string) (models.Rule, error) {
	var result models.Rule
	err := db.Table("rules").
		Where("name = ? AND alert = ?", name, alert).
		Scan(&result).Error
	return result, err
}

func (db *datastore) CreateAlertRule(rule *models.Rule) error {

	var result models.Rule
	notfound := db.Where("rules.name = ? AND rules.alert = ?", rule.Name, rule.Alert).
		First(&result).
		RecordNotFound()
	if !notfound {
		return errors.New("The rule is in Database")
	}

	return db.Save(rule).Error
}

func (db *datastore) UpdateAlertRule(rule *models.Rule) error {
	var result models.Rule
	notfound := db.Where("rules.name = ? AND rules.alert = ?", rule.Name, rule.Alert).
		First(&result).
		RecordNotFound()
	if notfound {
		return errors.New("The rule not found")
	} else {
		return db.Model(rule).
			Where("rules.name = ? AND rules.alert = ?", rule.Name, rule.Alert).
			Omit("name, alert").
			Updates(rule).Error
	}
}

func (db *datastore) DeleteAlertRuleByIDName(id uint64, name string) (int64, error) {
	result := db.Debug().Where("rules.id = ? AND rules.name = ?", id, name).
		Delete(&models.Rule{})
	err := result.Error
	rowsAffected := result.RowsAffected
	return rowsAffected, err
}
