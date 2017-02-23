package store

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type Store interface {
	ListAlertRules(page models.Page, name string) (map[string]interface{}, error)
	GetAlertRule(id uint64) (models.Rule, error)
	CreateAlertRule(rule *models.Rule) error
	UpdateAlertRule(rule *models.Rule) error
	DeleteAlertRule(id uint64) (int64, error)
	RuleNotFound(rule *models.Rule) bool
	CreateOrIncreaseEvent(event *models.Event) error
	AckEvent(pk int, username string, groupname string) error
	ListAckedEvent(page models.Page, username string, groupname string) map[string]interface{}
	ListUnackedEvent(page models.Page, username string, groupname string) map[string]interface{}
}
