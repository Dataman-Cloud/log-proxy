package prometheusrule

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/config"
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
	app := strings.Replace(rule.App, "-", "_", -1)
	alert := fmt.Sprintf("%s_%s_%s", rule.Class, rule.Name, app)
	pending := rule.Pending
	serverity := rule.Severity
	indicator := rule.Indicator
	aggregation := rule.Aggregation
	comparison := rule.Comparison
	threshold := strconv.FormatInt(rule.Threshold, 10)
	var unit string
	ruleIndicator, ok := ruleMap.mapper[rule.Indicator]
	if ok {
		switch indicator {
		case "mem_usage":
			unit = "%"
		case "cpu_usage":
			unit = "%"
		case "tomcat_thread_count":
			unit = ""
		}
	} else {
		return nil, errors.New("Cannot support monitor indicator: " + rule.Indicator)
	}
	judgement := fmt.Sprintf("%s %s %s%s", aggregation, comparison, threshold, unit)
	duration := rule.Duration
	labels := fmt.Sprintf(`{ cluster = "%s", app = "%s", value = "{{ $value }}", severity = "%s", indicator = "%s", judgement = "%s", duration = "%s" }`, rule.Cluster, rule.App, serverity, indicator, judgement, duration)
	annotations := `{ description = "", summary = "" }`

	raw := models.RawRule{}
	raw.Alert = alert
	raw.Pending = pending
	raw.Labels = labels
	raw.Annotations = annotations

	ruleIndicator, ok = ruleMap.mapper[rule.Indicator]
	if ok {
		templ := ruleIndicator.Templ
		switch rule.Indicator {
		case "mem_usage":
			raw.Expr = fmt.Sprintf(templ, aggregation, rule.Cluster, rule.App,
				rule.Cluster, rule.App, comparison, threshold)
		case "cpu_usage":
			raw.Expr = fmt.Sprintf(templ, aggregation, rule.Cluster, rule.App,
				rule.Duration, comparison, threshold)
		case "tomcat_thread_count":
			raw.Expr = fmt.Sprintf(templ, aggregation, rule.Cluster, rule.App,
				comparison, threshold)
		}
	} else {
		return nil, errors.New("Cannot support monitor indicator: " + rule.Indicator)
	}
	return &raw, nil
}

func (ruleMap *RuleMapper) GetRuleIndicatorsList() (keys map[string]string) {
	keys = make(map[string]string)
	for k, v := range GetRuleIndicators() {
		keys[k] = v.Unit
	}
	return keys
}

func Event2Cama(event *models.Event) *models.CamaEvent {
	var recover int
	if event.Ack == true {
		recover = 1
	} else {
		recover = 0
	}

	return &models.CamaEvent{
		ID:        event.AlertName,
		Channel:   "DOCKER",
		FirstTime: event.CreatedAt.Format(config.CamaTimeFormatString),
		LastTime:  event.UpdatedAt.Format(config.CamaTimeFormatString),
		Recover:   recover,
		Merger:    event.Count,
		Node:      "",
		NodeAlias: "",
		ServerNo:  "",
		EventDesc: event.Description,
		Level:     5,
	}
}
