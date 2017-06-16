package alertservice

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"
	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusrule"

	log "github.com/Sirupsen/logrus"
)

const (
	// ReceiveEventError code
	ReceiveEventError = "503-21000"
	// AckEventError code
	AckEventError = "503-21001"

	// PromtheusReloadPath path string
	PromtheusReloadPath = "/-/reload"

	ruleTempl = `# This rule was update at {{ .UpdatedAt }}
ALERT {{.Alert}}
  IF {{ .Expr }}
  FOR {{ .Pending }}
  LABELS {{ .Labels }}
  ANNOTATIONS {{ .Annotations }}
`
	ruleFileUpdate = "update"
	ruleFileDelete = "delete"
	ruleInterval   = "1m"

	ruleStatusActive   = "Enabled"
	ruleStatusInActive = "Disabled"
)

type Alert struct {
	Store      store.Store
	HTTPClient *http.Client
	PromServer string
	Interval   string
	RulesPath  string
	Rule       *models.Rule
	Indicators map[string]*models.Indicator
}

func initAlert() *Alert {
	var interval, promServer, rulesPath string
	if config.GetConfig() == nil {
		interval = ruleInterval
		promServer = ""
		rulesPath = ""
	} else {
		interval = config.GetConfig().RuleFileInterval
		if interval == "" {
			interval = ruleInterval
		}
		promServer = config.GetConfig().PrometheusURL
		rulesPath = config.GetConfig().RuleFilePath
	}
	return &Alert{
		Store:      datastore.From(database.GetDB()),
		HTTPClient: http.DefaultClient,
		PromServer: promServer,
		Interval:   interval,
		RulesPath:  rulesPath,
		Rule:       models.NewRule(),
		Indicators: prometheusrule.NewRuleMapper().GetRuleIndicators(),
	}
}

// NewAlert init the struct Alert
func NewAlert() service.Alerter {
	return initAlert()
}

// GetAlertIndicators return the indicator list
func (alert *Alert) GetAlertIndicators() map[string]string {
	var result map[string]string
	result = make(map[string]string)

	for alias, indicator := range alert.Indicators {
		result[alias] = indicator.Unit
	}

	return result
}

// CreateAlertRule create the alert rule in Database
func (alert *Alert) CreateAlertRule(rule *models.Rule) (*models.Rule, error) {
	var (
		result         models.Rule
		err            error
		app            string
		indicatorName  string
		indicatorAlias string
		indicatorUnit  string
	)

	if err = isValidRuleFile(rule); err != nil {
		return nil, err
	}
	// transfer the indicator from alias to name
	indicatorAlias = rule.Indicator
	indicatorName, indicatorUnit, err = alert.GetIndicatorName(indicatorAlias)
	if err != nil {
		return nil, err
	}
	rule.Indicator = indicatorName
	rule.Unit = indicatorUnit
	app = strings.Replace(rule.App, "-", "_", -1)
	rule.Name = fmt.Sprintf("%s_%s_%s", app, rule.Indicator, rule.Severity)
	// Create Alert rule in DB
	err = alert.Store.CreateAlertRule(rule)
	if err != nil {
		return nil, fmt.Errorf("Failed to create alert rule in DB with error %v", err)
	}
	// Get the rule record from DB
	result, err = alert.Store.GetAlertRuleByName(rule.Name)
	if err != nil {
		return nil, err
	}
	// Write the Rule in file and reload Prometheus conf
	err = alert.WriteAlertFile(rule)
	if err != nil {
		return nil, err
	}
	// Update the rule status as active
	result.Status = ruleStatusActive
	err = alert.Store.UpdateAlertRule(&result)
	if err != nil {
		return nil, err
	}

	// transfer the indicator from name to alias
	result.Indicator = indicatorAlias
	return &result, err
}

func (alert *Alert) GetIndicatorName(alias string) (name, unit string, err error) {
	if indicator, ok := alert.Indicators[alias]; ok {
		name = indicator.Name
		unit = indicator.Unit
		return name, unit, nil
	}
	return "", "", fmt.Errorf("The '%s' is not any of indicators's alias", alias)
}

func (alert *Alert) GetIndicatorAlias(name string) (alias, unit string, err error) {
	ruleMap := prometheusrule.NewRuleMapper()
	indicators := ruleMap.GetRuleIndicatorsInName()
	if indicator, ok := indicators[name]; ok {
		alias = indicator.Alias
		unit = indicator.Unit
		return alias, unit, nil
	}
	return "", "", fmt.Errorf("The '%s' is not any of indicators's name", name)
}

func isValidRuleFile(rule *models.Rule) error {
	switch {
	case rule.App == "":
		return fmt.Errorf("app required")
	case rule.Pending == "":
		return fmt.Errorf("pending required")
	case rule.Indicator == "":
		return fmt.Errorf("indicator required")
	case rule.Severity == "":
		return fmt.Errorf("severity required")
	case rule.Aggregation == "":
		return fmt.Errorf("aggregation required")
	case rule.Comparison == "":
		return fmt.Errorf("comparison required")
	}

	return nil
}

// ReloadPrometheusConf reload the conf by calling prometheus api
func (alert *Alert) ReloadPrometheusConf() error {
	u, err := url.Parse(alert.PromServer)
	if err != nil {
		return err
	}
	u.Path = strings.TrimRight(u.Path, "/") + PromtheusReloadPath
	resp, err := alert.HTTPClient.Post(u.String(), "application/json", nil)
	if err != nil || resp.StatusCode != 200 {
		err = fmt.Errorf("Failed to reload the configuration file of prometheus %s, return %d", u.String(), resp.StatusCode)
		return err
	}
	defer resp.Body.Close()

	return nil
}

// WriteAlertFile write the alert rule to file
func (alert *Alert) WriteAlertFile(rule *models.Rule) error {
	var (
		mapper  *prometheusrule.RuleMapper
		rawRule *models.RawRule
		err     error
	)
	// mapping the rule from rule to raw rule
	mapper = prometheusrule.NewRuleMapper()
	rawRule, err = mapper.Map2Raw(rule)
	if err != nil {
		return err
	}
	rawRule.UpdatedAt = time.Now()

	//open the alert rule file
	path := alert.RulesPath
	alertfile := fmt.Sprintf("%s/%s.rule", path, rawRule.Alert)
	f, err := os.Create(alertfile)
	defer f.Close()
	if err != nil {
		return err
	}
	// convert the rawRule with the template
	t := template.Must(template.New("ruleTempl").Parse(ruleTempl))
	var buf bytes.Buffer
	err = t.Execute(&buf, rawRule)
	if err != nil {
		log.Errorln("executing templta: ", err)
		return err
	}

	f.WriteString(buf.String())

	err = alert.ReloadPrometheusConf()
	if err != nil {
		return err
	}

	return nil
}

func (alert *Alert) DeleteAlertRule(id uint64, group string) error {
	var (
		rowsAffected int64
		err          error
		result       models.Rule
	)

	// Get alert rule by ID
	result, err = alert.Store.GetAlertRule(id)
	if err != nil {
		return fmt.Errorf("DeleteAlertRule: GetAlertRule() %v, ", err)
	}

	if result.Group != group {
		return fmt.Errorf("DeleteAlertRule: Can't delete the rule %d with group %s", id, group)
	}

	// Delate alert rule
	rowsAffected, err = alert.Store.DeleteAlertRuleByID(id)
	if err != nil {
		return fmt.Errorf("Failed to delete the rule by %v", err)
	}

	if rowsAffected == 0 {
		return nil
	}

	// Update the alert file content
	err = alert.UpdateAlertFile(&result)
	if err != nil {
		return fmt.Errorf("DeleteAlertRule: Delete Alert file error %v", err)
	}

	return err
}

// UpdateAlertFile remove the rule from the file.
func (alert *Alert) UpdateAlertFile(rule *models.Rule) error {
	var (
		err     error
		message string
	)

	filename := rule.Name

	path := alert.RulesPath
	alertfile := fmt.Sprintf("%s/%s.rule", path, filename)
	f, err := os.Create(alertfile)
	defer f.Close()
	if err != nil {
		return err
	}

	message = "# inactive this rule"

	f.WriteString(message + "\n")

	err = alert.ReloadPrometheusConf()
	if err != nil {
		return err
	}
	return nil
}

// ListAlertRules list the rules by name with pages.
func (alert *Alert) ListAlertRules(page models.Page, groups []string, app string) (*models.RulesList, error) {
	var (
		result         *models.RulesList
		err            error
		indicatorName  string
		indicatorAlias string
	)
	result, err = alert.Store.ListAlertRules(page, groups, app)
	if err != nil {
		return nil, err
	}
	for _, rule := range result.Rules {
		indicatorName = rule.Indicator
		indicatorAlias, _, err = alert.GetIndicatorAlias(indicatorName)
		if err != nil {
			return nil, err
		}
		rule.Indicator = indicatorAlias
	}
	return result, err
}

// GetAlertRule return the info of alert rule by id
func (alert *Alert) GetAlertRule(id uint64) (*models.Rule, error) {
	var (
		result         models.Rule
		rule           *models.Rule
		err            error
		indicatorName  string
		indicatorAlias string
	)

	result, err = alert.Store.GetAlertRule(id)
	if err != nil {
		return &result, err
	}

	rule = &result
	indicatorName = rule.Indicator
	indicatorAlias, _, err = alert.GetIndicatorAlias(indicatorName)
	if err != nil {
		return nil, err
	}
	rule.Indicator = indicatorAlias

	return rule, err
}

// UpdateAlertRule update the alert rule in Database
func (alert *Alert) UpdateAlertRule(rule *models.Rule) (*models.Rule, error) {
	var (
		result         models.Rule
		indicatorName  string
		indicatorAlias string
		err            error
	)

	id := rule.ID
	group := rule.Group

	// Can't Update the app, severity and indicator of one rule, if required, delete and create it.
	if rule.App != "" || rule.Severity != "" || rule.Indicator != "" {
		return nil, fmt.Errorf("Don't Allow to update the app, servirity and indicator of the rule")
	}

	// Get alert rule by ID
	result, err = alert.Store.GetAlertRule(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get the rule %d with error %v", id, err)
	}

	if result.Group != group {
		return nil, fmt.Errorf("Don't Allow to update the rule %d with group %s, it is in the group %s", id, group, result.Group)
	}

	err = alert.Store.UpdateAlertRule(rule)
	if err != nil {
		return nil, fmt.Errorf("Failed to update the rule %d", id)
	}

	result, err = alert.Store.GetAlertRule(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get the updated rule %d with error %v. The prometheus conf isn't updated", id, err)
	}
	if result.Status == ruleStatusInActive {
		err = alert.UpdateAlertFile(&result)
		if err != nil {
			return nil, fmt.Errorf("Failed to inactive the alert conf in prometheus with error %v", err)
		}
	} else {
		err = alert.WriteAlertFile(&result)
		if err != nil {
			return nil, fmt.Errorf("Failed to update the alert conf in prometheus with error %v", err)

		}
	}
	rule = &result
	indicatorName = rule.Indicator
	indicatorAlias, _, err = alert.GetIndicatorAlias(indicatorName)
	if err != nil {
		return nil, err
	}
	rule.Indicator = indicatorAlias

	return rule, err
}
