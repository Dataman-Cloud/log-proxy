package service

import "github.com/Dataman-Cloud/log-proxy/src/models"

type Alerter interface {
	GetAlertIndicators() map[string]string
	CreateAlertRule(rule *models.Rule) (*models.Rule, error)
	GetIndicatorName(alias string) (name, unit string, err error)
	GetIndicatorAlias(name string) (alias, unit string, err error)
	ReloadPrometheusConf() error
	WriteAlertFile(rule *models.Rule) error
	DeleteAlertRule(id uint64, group string) error
	UpdateAlertFile(rule *models.Rule) error
	ListAlertRules(page models.Page, group []string, app string) (*models.RulesList, error)
	GetAlertRule(id uint64) (*models.Rule, error)
	UpdateAlertRule(rule *models.Rule) (*models.Rule, error)
	ReceiveAlertEvent(message map[string]interface{}) error
	GetAlertEvents(page models.Page, options map[string]interface{}) (map[string]interface{}, error)
	AckAlertEvent(id int, options map[string]interface{}) error
}
