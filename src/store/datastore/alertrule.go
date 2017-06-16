package datastore

import (
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) ListAlertRules(page models.Page, opts map[string]interface{}) (*models.RulesList, error) {
	var (
		rules     []*models.Rule
		count     int
		err       error
		rulesList *models.RulesList
	)
	db.Table("rules").
		Where(opts).
		Scan(&rules).
		Count(&count)
	err = db.Debug().Table("rules").
		Where(opts).
		Offset(page.PageFrom).
		Limit(page.PageSize).
		Scan(&rules).
		Error
	rulesList = models.NewRulesList()
	rulesList.Count = int64(count)
	rulesList.Rules = rules

	return rulesList, err
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

func (db *datastore) GetAlertRuleByName(name string) (models.Rule, error) {
	var result models.Rule
	err := db.Table("rules").
		Where("name = ?", name).
		Scan(&result).Error
	return result, err
}

func (db *datastore) CreateAlertRule(rule *models.Rule) error {
	var result models.Rule
	notfound := db.Where("rules.Name = ?", rule.Name).
		First(&result).
		RecordNotFound()
	if !notfound {
		return errors.New("The rule is in Database")
	}

	return db.Save(rule).Error
}

func (db *datastore) UpdateAlertRule(rule *models.Rule) error {
	var result *models.Rule
	notfound := db.Where("rules.id = ?", rule.ID).
		First(&result).
		RecordNotFound()
	if notfound {
		return errors.New("The rule not found")
	}
	return db.Model(rule).
		Where("rules.id = ?", rule.ID).
		Omit("name, app, severity, indicator").
		Updates(rule).Error
}

func (db *datastore) DeleteAlertRuleByID(id uint64) (int64, error) {
	result := db.Debug().Where("rules.id = ?", id).
		Delete(&models.Rule{})
	err := result.Error
	rowsAffected := result.RowsAffected
	return rowsAffected, err
}
