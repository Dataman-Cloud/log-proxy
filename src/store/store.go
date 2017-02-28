package store

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type Store interface {
	ListAlertRules(page models.Page, user string) (map[string]interface{}, error)
	GetAlertRule(id uint64) (models.Rule, error)
	GetAlertRules() ([]*models.Rule, error)
	GetAlertRuleByName(name, alert string) (models.Rule, error)
	CreateAlertRule(rule *models.Rule) error
	UpdateAlertRule(rule *models.Rule) error
	DeleteAlertRuleByIDName(id uint64, name string) (int64, error)
	CreateOrIncreaseEvent(event *models.Event) error
	AckEvent(pk int, username string, groupname string) error
	ListAckedEvent(page models.Page, username string, groupname string) map[string]interface{}
	ListUnackedEvent(page models.Page, username string, groupname string) map[string]interface{}
}
