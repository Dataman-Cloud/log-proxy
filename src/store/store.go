package store

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type Store interface {
	ListAlertRules() ([]*models.Rule, error)
	GetAlertRule(id uint64) (models.Rule, error)
	GetAlertRuleByName(name, alert string) (models.Rule, error)
	CreateAlertRule(rule *models.Rule) error
	UpdateAlertRule(rule *models.Rule) error
	DeleteAlertRule(id uint64) (int64, error)
	DeleteAlertRuleByName(name, alert string) (int64, error)
	ValidataRule(rule *models.Rule) error
}
