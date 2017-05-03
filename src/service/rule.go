package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"
	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusrule"
	"github.com/prometheus/common/log"
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

// NewAlert init the struct Alert
func NewAlert() *Alert {
	interval := config.GetConfig().RuleFileInterval
	if interval == "" {
		interval = ruleInterval
	}

	return &Alert{
		Store:      datastore.From(database.GetDB()),
		HTTPClient: http.DefaultClient,
		PromServer: config.GetConfig().PrometheusURL,
		Interval:   interval,
		RulesPath:  config.GetConfig().RuleFilePath,
		Rule:       models.NewRule(),
		Indicators: prometheusrule.NewRuleMapper().GetRuleIndicators(),
	}
}

// CreateAlertRule create the alert rule in Database
func (alert *Alert) CreateAlertRule(rule *models.Rule) (*models.Rule, error) {
	var (
		result                             models.Rule
		err                                error
		app, indicatorName, indicatorAlias string
	)

	if err = isValidRuleFile(rule); err != nil {
		return nil, err
	}
	// transfer the indicator from alias to name
	indicatorAlias = rule.Indicator
	indicatorName, err = alert.getIndicatorName(indicatorAlias)
	if err != nil {
		return nil, err
	}
	rule.Indicator = indicatorName

	app = strings.Replace(rule.App, "-", "_", -1)
	rule.Name = fmt.Sprintf("%s_%s_%s", app, rule.Indicator, rule.Severity)

	// Create Alert rule in DB
	err = alert.Store.CreateAlertRule(rule)
	if err != nil {
		return nil, err
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

func (alert *Alert) getIndicatorName(alias string) (name string, err error) {
	if indicator, ok := alert.Indicators[alias]; ok {
		name = indicator.Name
		return name, nil
	}
	return "", fmt.Errorf("The '%s' is not any of indicators's alias", alias)
}

func (alert *Alert) getIndicatorAlias(name string) (alias string, err error) {
	ruleMap := prometheusrule.NewRuleMapper()
	indicators := ruleMap.GetRuleIndicatorsInName()
	if indicator, ok := indicators[name]; ok {
		alias = indicator.Alias
		return alias, nil
	}
	return "", fmt.Errorf("The '%s' is not any of indicators's name", name)
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
