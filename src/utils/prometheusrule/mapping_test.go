package prometheusrule

import (
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func TestNewRuleMapper(t *testing.T) {
	data := NewRuleMapper()
	if data == nil {
		t.Errorf("Expect data is not nil, got %v", data)
	}
}

func TestGetRuleIndicators(t *testing.T) {
	ruleMap := NewRuleMapper()
	data := ruleMap.GetRuleIndicators()
	if len(data) != 3 {
		t.Errorf("expect 3, but got %d", len(data))
	}
}

func TestMap2Raw(t *testing.T) {
	rule := models.NewRule()
	rule.Name = "work_nginx_cpu_usage_warning"
	rule.App = "work-nginx"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	ruleMap := NewRuleMapper()
	rawRule, err := ruleMap.Map2Raw(rule)
	if rawRule.Alert != "work_nginx_cpu_usage_warning" {
		t.Errorf("expect the raw rule alert work_nginx_cpu_usage_warning, but got %s", rawRule.Alert)
	}
	if err != nil {
		t.Errorf("expect the err is nil, but got %v", err)
	}
}
