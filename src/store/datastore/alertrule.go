package datastore

import (
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) ListAlertRules(page models.Page, name string) (map[string]interface{}, error) {
	var (
		rules []*models.Rule
		count int
	)
	err := db.Table("rules").Debug().Where("name = ?", name).Find(&rules).Count(&count).Error
	err = db.Table("rules").Debug().Where("name = ?", name).Offset(page.PageFrom).Limit(page.PageSize).Find(&rules).Error

	return map[string]interface{}{"count": count, "rules": rules}, err
}

func (db *datastore) GetAlertRule(id uint64) (models.Rule, error) {
	var result models.Rule
	fields := []string{"id", "updated_at", "deleted_at", "name", "alert", "expr", "duration", "labels", "description", "summary"}
	err := db.Table("rules").
		Select(fields).Where("id = ?", id).
		Scan(&result).Error
	return result, err
}

func (db *datastore) GetAlertRuleByName(name, alert string) (models.Rule, error) {
	var result models.Rule
	fields := []string{"id", "updated_at", "deleted_at", "name", "alert", "expr", "duration", "labels", "description", "summary"}
	err := db.Table("rules").
		Select(fields).Where("name = ? AND alert = ?", name, alert).
		Scan(&result).Error
	return result, err
}

func (db *datastore) CreateAlertRule(rule *models.Rule) error {
	return db.Save(rule).Error
}

func (db *datastore) UpdateAlertRule(rule *models.Rule) error {
	return db.Model(rule).
		Where("name = ? AND alert = ?", rule.Name, rule.Alert).
		Omit("name", "alert").
		Updates(rule).Error
}

func (db *datastore) DeleteAlertRule(id uint64) (int64, error) {
	result := db.Where("rules.id = ?", id).
		Delete(&models.Rule{})
	err := result.Error
	rowsAffected := result.RowsAffected
	return rowsAffected, err
}

func (db *datastore) DeleteAlertRuleByName(name, alert string) (int64, error) {
	result := db.Where("rules.name = ? AND rules.alert", name, alert).
		Delete(&models.Rule{})
	err := result.Error
	rowsAffected := result.RowsAffected
	return rowsAffected, err
}

func (db *datastore) ValidataRule(rule *models.Rule) error {
	var result models.Rule
	result, _ = db.GetAlertRuleByName(rule.Name, rule.Alert)
	fmt.Printf("result: %s", result)

	if result.Alert == rule.Alert && result.Name == rule.Name {
		return nil
	}
	err := fmt.Errorf("Rule name %s and alert %s is not validata.", rule.Name, rule.Alert)
	return err
}
