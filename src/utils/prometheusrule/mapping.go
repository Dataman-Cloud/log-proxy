package prometheusrule

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type RuleMapper struct {
	appRegex       *regexp.Regexp
	severityRegex  *regexp.Regexp
	indicatorRegex *regexp.Regexp
	judgementRegex *regexp.Regexp
	durationRegex  *regexp.Regexp
}

func NewRuleMapper() *RuleMapper {
	return &RuleMapper{
		appRegex:       regexp.MustCompile(`container_label_APP="(.*?)"`),
		severityRegex:  regexp.MustCompile(`severity="(.*?)"`),
		indicatorRegex: regexp.MustCompile(`indicator="(.*?)"`),
		judgementRegex: regexp.MustCompile(`judgement="(.*?)"`),
		durationRegex:  regexp.MustCompile(`duration="(.*?)"`),
	}
}

// Map2Raw convert the rule to rawRule
func (ruleMap *RuleMapper) Map2Raw(rule *models.Rule) (*models.RawRule, error) {
	alert := rule.Name
	pending := rule.Pending
	serverity := rule.Severity
	indicator := rule.Indicator
	aggregation := rule.Aggregation
	comparison := rule.Comparison
	threshold := strconv.FormatInt(rule.Threshold, 10)
	var unit string
	mapper := ruleMap.GetRuleIndicatorsInName()
	ruleIndicator, ok := mapper[rule.Indicator]
	if ok {
		unit = ruleIndicator.Unit
	} else {
		return nil, errors.New("Cannot support monitor indicator: " + rule.Indicator)
	}
	judgement := fmt.Sprintf("%s %s %s%s", aggregation, comparison, threshold, unit)
	duration := rule.Duration
	labels := fmt.Sprintf(`{ app = "%s", value = "{{ $value }}", severity = "%s", indicator = "%s", judgement = "%s", duration = "%s" }`, rule.App, serverity, indicator, judgement, duration)
	annotations := `{ description = "", summary = "" }`

	raw := models.RawRule{}
	raw.Alert = alert
	raw.Pending = pending
	raw.Labels = labels
	raw.Annotations = annotations

	ruleIndicator, ok = mapper[rule.Indicator]
	if ok {
		templ := ruleIndicator.Templ
		switch rule.Indicator {
		case "mem_usage":
			raw.Expr = fmt.Sprintf(templ, aggregation, rule.App,
				rule.App, comparison, threshold)
		case "cpu_usage":
			raw.Expr = fmt.Sprintf(templ, aggregation, rule.App,
				rule.Duration, comparison, threshold)
		case "tomcat_thread_count":
			raw.Expr = fmt.Sprintf(templ, aggregation, rule.App,
				comparison, threshold)
		}
	} else {
		return nil, errors.New("Cannot support monitor indicator: " + rule.Indicator)
	}
	return &raw, nil
}

func (ruleMap *RuleMapper) GetRuleIndicators() map[string]*models.Indicator {
	var ruleIndicators = make(map[string]*models.Indicator)
	ruleIndicators[cpuUsagePercent.Alias] = cpuUsagePercent
	ruleIndicators[memUsagePercent.Alias] = memUsagePercent
	ruleIndicators[tomcatThreadPool.Alias] = tomcatThreadPool
	return ruleIndicators
}

func (ruleMap *RuleMapper) GetRuleIndicatorsInName() map[string]*models.Indicator {
	var ruleIndicators = make(map[string]*models.Indicator)
	ruleIndicators[cpuUsagePercent.Name] = cpuUsagePercent
	ruleIndicators[memUsagePercent.Name] = memUsagePercent
	ruleIndicators[tomcatThreadPool.Name] = tomcatThreadPool
	return ruleIndicators
}
