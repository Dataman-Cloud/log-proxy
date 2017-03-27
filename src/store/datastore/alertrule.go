package datastore

import (
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

const (
	queryIndex = "rules.class = ? AND rules.name = ? AND rules.cluster = ? AND rules.app = ?"
)

func (db *datastore) CreateAlertRule(rule *models.Rule) error {
	var result models.Rule
	notfound := db.Where(queryIndex, rule.Class, rule.Name, rule.Cluster, rule.App).
		First(&result).
		RecordNotFound()
	if !notfound {
		return fmt.Errorf("This rule is in Database")
	}
	return db.Save(rule).Error
}

func (db *datastore) GetAlertRuleByUniqueIndex(class, name, cluster, app string) (models.Rule, error) {
	var result models.Rule
	err := db.Table("rules").
		Where(queryIndex, class, name, cluster, app).
		Scan(&result).Error
	return result, err
}

func (db *datastore) ListAlertRules(page models.Page, class, cluster, app string) (map[string]interface{}, error) {
	var (
		rules []*models.Rule
		count int
		err   error
	)
	if class == "" {
		db.Table("rules").Find(&rules).Count(&count)
		err = db.Table("rules").Offset(page.PageFrom).Limit(page.PageSize).Find(&rules).Error
	} else if class != "" && cluster != "" && app == "" {
		db.Table("rules").
			Where("class = ? AND cluster = ?", class, cluster).
			Scan(&rules).
			Count(&count)
		err = db.Table("rules").
			Where("class = ? AND cluster = ?", class, cluster).
			Offset(page.PageFrom).
			Limit(page.PageSize).
			Scan(&rules).
			Error
	} else if class != "" && cluster != "" && app != "" {
		db.Table("rules").
			Where("class = ? AND cluster = ? AND app = ?", class, cluster, app).
			Scan(&rules).
			Count(&count)
		err = db.Table("rules").
			Where("class = ? AND cluster = ? AND app = ?", class, cluster, app).
			Offset(page.PageFrom).
			Limit(page.PageSize).
			Scan(&rules).
			Error
	} else {
		rules = nil
	}
	return map[string]interface{}{"count": count, "rules": rules}, err
}

func (db *datastore) GetAlertRule(id uint64) (models.Rule, error) {
	var result models.Rule
	err := db.Table("rules").
		Where("id = ?", id).
		Scan(&result).Error
	return result, err
}

func (db *datastore) DeleteAlertRuleByIDClass(id uint64, class string) (int64, error) {
	result := db.Where("rules.id = ? AND rules.class = ?", id, class).
		Delete(&models.Rule{})
	err := result.Error
	rowsAffected := result.RowsAffected
	return rowsAffected, err
}

func (db *datastore) UpdateAlertRule(rule *models.Rule) error {
	var result models.Rule
	notfound := db.
		Where("id = ?", rule.ID).
		First(&result).
		RecordNotFound()
	if notfound {
		return fmt.Errorf("This rule is not in Database")
	}

	notfound = db.
		Where(queryIndex, rule.Class, rule.Name, rule.Cluster, rule.App).
		First(&result).
		RecordNotFound()
	if notfound {
		return fmt.Errorf("Can't update the rule fields class, name, cluster and app")
	}

	return db.Model(rule).
		Where("id = ?", rule.ID).
		Omit("class, name, cluster, app").
		Updates(rule).Error
}

func (db *datastore) GetAlertRules() ([]*models.Rule, error) {
	var (
		rules []*models.Rule
		err   error
	)
	err = db.Table("rules").Find(&rules).Error

	return rules, err
}
