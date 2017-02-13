package service

import (
	"fmt"
	"strconv"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/jinzhu/gorm"
)

type AlertRule struct {
	DB   *gorm.DB
	Rule *models.Rule
}

type Result struct {
	Name  string
	Alert string
}

func (AlertRule *AlertRule) CreateAlertRule() error {
	db := AlertRule.DB
	rule := AlertRule.Rule

	var result Result
	db.Table("rules").Select("name, alert").Where("name = ? AND alert = ?", rule.Name, rule.Alert).Scan(&result)

	if result.Alert == rule.Alert && result.Name == rule.Name {
		err := fmt.Errorf("Rule name %s and alert %s is exsit.", rule.Name, rule.Alert)
		return err
	}

	err := db.Create(rule).Error

	return err
}

func (AlertRule *AlertRule) DeleteAlertRule(id string) error {
	db := AlertRule.DB

	if ruleID, err := strconv.ParseUint(id, 10, 32); err == nil {
		db.Where("id = ?", ruleID).Delete(&models.Rule{})
	} else {
		return err
	}
	return nil
}

func (AlertRule *AlertRule) GetAlertRule() ([]*models.Rule, error) {
	db := AlertRule.DB

	var rules []*models.Rule
	fields := []string{"id", "updated_at", "name", "alert", "for_time", "labels", "description", "summary"}
	err := db.Table("rules").Select(fields).Where("deleted_at IS NULL").Scan(&rules).Error

	return rules, err
}

func (AlertRule *AlertRule) UpdateAlertRule() error {
	db := AlertRule.DB
	rule := AlertRule.Rule

	var result Result
	db.Table("rules").Select("name, alert").Where("name = ? AND alert = ?", rule.Name, rule.Alert).Scan(&result)

	if result.Alert == rule.Alert && result.Name == rule.Name {
		err := db.Model(rule).Where("name = ? AND alert = ?", rule.Name, rule.Alert).Omit("name", "alert").Updates(rule).Error
		return err
	}

	err := fmt.Errorf("Rule name %s and alert %s is not exsit.", rule.Name, rule.Alert)
	return err
}
