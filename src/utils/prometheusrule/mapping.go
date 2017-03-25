package prometheusrule

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type RuleMapper struct {
	mapper         map[string]*models.Indicator
	clusterRegex   *regexp.Regexp
	appRegex       *regexp.Regexp
	severityRegex  *regexp.Regexp
	indicatorRegex *regexp.Regexp
	judgementRegex *regexp.Regexp
	durationRegex  *regexp.Regexp
}

func NewRuleMapper() *RuleMapper {
	return &RuleMapper{
		mapper:         GetRuleIndicators(),
		clusterRegex:   regexp.MustCompile(`container_label_VCLUSTER="(.*?)"`),
		appRegex:       regexp.MustCompile(`container_label_APP="(.*?)"`),
		severityRegex:  regexp.MustCompile(`severity="(.*?)"`),
		indicatorRegex: regexp.MustCompile(`indicator="(.*?)"`),
		judgementRegex: regexp.MustCompile(`judgement="(.*?)"`),
		durationRegex:  regexp.MustCompile(`duration="(.*?)"`),
	}
}

func (ruleMap *RuleMapper) Map2Raw(rule *models.Rule) (*models.RawRule, error) {
	var app string
	if appName := strings.Split(rule.App, "-"); len(appName) == 2 {
		app = appName[1]
	} else {
		app = rule.App
	}
	alert := fmt.Sprintf("%s_%s_%s_%s", rule.Class, rule.Name, rule.Cluster, app)
	pending := rule.Pending
	serverity := rule.Severity
	indicator := rule.Indicator
	aggregation := rule.Aggregation
	comparison := rule.Comparison
	threshold := strconv.FormatInt(rule.Threshold, 10)
	judgement := fmt.Sprintf("%s %s %s", aggregation, comparison, threshold)
	duration := rule.Duration
	labels := fmt.Sprintf(`{ value = "{{ $value }}", severity = "%s", indicator = "%s", judgement = "%s", duration = "%s" }`, serverity, indicator, judgement, duration)
	annotations := `{ description = "", summary = "" }`

	raw := models.RawRule{}
	raw.Alert = alert
	raw.Pending = pending
	raw.Labels = labels
	raw.Annotations = annotations

	ruleIndicator, ok := ruleMap.mapper[rule.Indicator]
	if ok {
		templ := ruleIndicator.Templ
		switch rule.Indicator {
		case "mem_usage":
			raw.Expr = fmt.Sprintf(templ, aggregation, rule.Cluster, rule.App,
				rule.Cluster, rule.App, comparison, threshold)
		case "cpu_usage":
			raw.Expr = fmt.Sprintf(templ, aggregation, rule.Cluster, rule.App,
				rule.Duration, comparison, threshold)
		}
	} else {
		return nil, errors.New("Cannot support monitor indicator: " + rule.Indicator)
	}
	return &raw, nil
}

func (ruleMap *RuleMapper) GetRuleIndicatorsList() (keys []string) {
	for k := range GetRuleIndicators() {
		keys = append(keys, k)
	}
	return keys
}
