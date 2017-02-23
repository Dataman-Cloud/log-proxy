package datastore

import (
	"errors"
	"fmt"

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
	if name == "" {
		err = db.Table("rules").
			Where("deleted_at IS NULL").
			Scan(&rules).Error
	} else {
		err = db.Table("rules").
			Where("name = ? AND deleted_at IS NULL", name).
			Scan(&rules).Error
	}
>>>>>>> Fix the function of validata/delete rules

	return map[string]interface{}{"count": count, "rules": rules}, err
}

func (db *datastore) GetAlertRule(id uint64) (models.Rule, error) {
	var result models.Rule
	err := db.Table("rules").
		Where("id = ?", id).
		Scan(&result).Error
	return result, err
}

func (db *datastore) CreateAlertRule(rule *models.Rule) error {
	return db.Save(rule).Error
}

func (db *datastore) UpdateAlertRule(rule *models.Rule) error {
	return db.Model(rule).
		Where("name = ?", rule.Name).
		Omit("name").
		Updates(rule).Error
}

func (db *datastore) DeleteAlertRule(id uint64) (int64, error) {
	var rule models.Rule
	recordNotFound := db.Debug().
		Where("rules.id = ?", id).
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

func (db *datastore) RuleNotFound(rule *models.Rule) bool {
	var result models.Rule
	fmt.Println("name: ", rule.Name)
	recordNotFound := db.Where("rules.name = ?", rule.Name).
		First(&result).
		RecordNotFound()
	fmt.Println("Notfound: ", recordNotFound)
	return recordNotFound
}
