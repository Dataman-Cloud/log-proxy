package store

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type Store interface {
	CreateAlertRule(rule *models.Rule) error
	GetAlertRuleByUniqueIndex(group, name, cluster, app string) (models.Rule, error)
	ListAlertRules(page models.Page, class, cluster, app string) (map[string]interface{}, error)
	GetAlertRule(id uint64) (models.Rule, error)
	GetAlertRules() ([]*models.Rule, error)
	UpdateAlertRule(rule *models.Rule) error
	DeleteAlertRuleByIDClass(id uint64, class string) (int64, error)
	CreateOrIncreaseEvent(event *models.Event) error
	AckEvent(pk int, username string, groupname string) error
	ListEvents(page models.Page, options map[string]interface{}) (map[string]interface{}, error)
	CreateLogAlertRule(rule *models.LogAlertRule) error
	UpdateLogAlertRule(rule *models.LogAlertRule) error
	DeleteLogAlertRule(ID string) error
	GetLogAlertRule(ID string) (models.LogAlertRule, error)
	GetLogAlertRules(page models.Page) (map[string]interface{}, error)
	CreateLogAlertEvent(event *models.LogAlertEvent) error
	GetLogAlertEvents(options map[string]interface{}, page models.Page) (map[string]interface{}, error)
	DeleteLogAlertEvents(start, end string) error
	GetLogAlertClusters(start, end string) ([]*models.LogAlertClusters, error)
	GetLogAlertApps(cluster, start, end string) ([]*models.LogAlertApps, error)
	AckLogAlertEvent(ID string) error
}
