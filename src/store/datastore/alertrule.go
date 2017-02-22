package datastore

import (
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

<<<<<<< HEAD
func (db *datastore) ListAlertRules(page models.Page, name string) (map[string]interface{}, error) {
	var (
		rules []*models.Rule
		count int
	)
	err := db.Table("rules").Debug().Where("name = ?", name).Find(&rules).Count(&count).Error
	err = db.Table("rules").Debug().Where("name = ?", name).Offset(page.PageFrom).Limit(page.PageSize).Find(&rules).Error
=======
func (db *datastore) ListAlertRules(name string) ([]*models.Rule, error) {
	var rules []*models.Rule
	var err error
	fields := []string{"id", "updated_at", "deleted_at", "name", "alert", "expr", "duration", "labels", "description", "summary"}
	if name == "" {
		err = db.Table("rules").
			Select(fields).
			Where("deleted_at IS NULL").
			Scan(&rules).Error
	} else {
		err = db.Table("rules").
			Select(fields).
			Where("name = ? AND deleted_at IS NULL", name).
			Scan(&rules).Error
	}
>>>>>>> Fix the function of validata/delete rules

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
	var rule models.Rule
	recordNotFound := db.
		Where("rules.id = ? AND deleted_at IS NULL", id).
		First(&rule).
		RecordNotFound()

	if recordNotFound {
		return 0, errors.New("No this rule in database")
	}

	result := db.Table("rules").
		Where("rules.id = ?", id).
		Delete(&models.Rule{})
	err := result.Error
	rowsAffected := result.RowsAffected
	return rowsAffected, err
}

func (db *datastore) DeleteAlertRuleByName(name, alert string) (int64, error) {
	recordNotFound := db.Table("rules").
		Where("rules.name = ? AND rules.alert AND deleted_at IS NULL", name, alert).
		RecordNotFound()
	if recordNotFound {
		return 0, errors.New("No this rule in database")
	}

	result := db.Where("rules.name = ? AND rules.alert", name, alert).
		Delete(&models.Rule{})
	err := result.Error
	rowsAffected := result.RowsAffected
	return rowsAffected, err
}

func (db *datastore) ValidataRule(rule *models.Rule) bool {
	var result models.Rule
	recordNotFound := db.Where("rules.name = ? AND rules.alert = ? AND deleted_at IS NULL", rule.Name, rule.Alert).
		First(&result).
		RecordNotFound()
	return !recordNotFound //if the rule in DB, return true
}
