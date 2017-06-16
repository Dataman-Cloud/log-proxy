package store

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type Store interface {
	ListAlertRules(page models.Page, groups []string, app string) (*models.RulesList, error)
	GetAlertRule(id uint64) (models.Rule, error)
	GetAlertRules() ([]*models.Rule, error)
	GetAlertRuleByName(name string) (models.Rule, error)
	CreateAlertRule(rule *models.Rule) error
	UpdateAlertRule(rule *models.Rule) error
	DeleteAlertRuleByID(id uint64) (int64, error)
	CreateOrIncreaseEvent(event *models.Event) error
	AckEvent(ID int, group, app string) error
	ListEvents(page models.Page, options map[string]interface{}, groups []string) (map[string]interface{}, error)
	CreateLogAlertRule(rule *models.LogAlertRule) error
	UpdateLogAlertRule(rule *models.LogAlertRule) error
	DeleteLogAlertRule(ID string) error
	GetLogAlertRule(ID string) (models.LogAlertRule, error)
	GetLogAlertRules(opts map[string]interface{}, page models.Page) (map[string]interface{}, error)
	GetLogAlertApps(opts map[string]interface{}, page models.Page) ([]*models.LogAlertApps, error)
	AckLogAlertEvent(ID string) error
	GetLogAlertEvent(ID string) (*models.LogAlertEvent, error)
	GetLogAlertEvents(opts map[string]interface{}, page models.Page) (map[string]interface{}, error)
	CreateLogAlertEvent(event *models.LogAlertEvent) error
}
