package prometheusrule

import (
	"fmt"
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

	rule.Name = "work_nginx_mem_usage_warning"
	rule.Indicator = "mem_usage"
	rawRule, err = ruleMap.Map2Raw(rule)
	if rawRule.Alert != "work_nginx_mem_usage_warning" {
		t.Errorf("expect the raw rule alert work_nginx_mem_usage_warning, but got %s", rawRule.Alert)
	}
	if err != nil {
		t.Errorf("expect the err is nil, but got %v", err)
	}

	rule.Name = "work_nginx_tomcat_thread_count_warning"
	rule.Indicator = "tomcat_thread_count"
	rawRule, err = ruleMap.Map2Raw(rule)
	if rawRule.Alert != "work_nginx_tomcat_thread_count_warning" {
		t.Errorf("expect the raw rule alert work_nginx_mem_usage_warning, but got %s", rawRule.Alert)
	}
	fmt.Println(rawRule.Expr)
	if err != nil {
		t.Errorf("expect the err is nil, but got %v", err)
	}
}

func TestRuleIndicator(t *testing.T) {
	var name, expectVaule string
	var indicator *models.Indicator
	name = "cpu_usage"
	indicator = RuleIndicator(name)
	expectVaule = `%s(irate(container_cpu_usage_seconds_total{container_label_DM_APP_ID="%s", container_label_DM_LOG_TAG!="ignore", id=~"/docker/.*", name=~"mesos.*"}[%s])) by( container_label_DM_APP_ID, container_label_DM_SLOT_INDEX ) keep_common * 100 %s %s`
	if indicator.Templ != expectVaule {
		t.Errorf("expect the templ is %s\n, got %s", expectVaule, indicator.Templ)
	}
	name = "mem_usage"
	indicator = RuleIndicator(name)
	expectVaule = `%s(container_memory_usage_bytes{container_label_DM_APP_ID="%s", container_label_DM_LOG_TAG!="ignore", id=~"/docker/.*", name=~"mesos.*"} / container_spec_memory_limit_bytes{container_label_DM_APP_ID="%s", container_label_DM_LOG_TAG!="ignore", id=~"/docker/.*", name=~"mesos.*"}) by( container_label_DM_APP_ID, container_label_DM_SLOT_INDEX ) keep_common * 100 %s %s`
	if indicator.Templ != expectVaule {
		t.Errorf("expect the templ is %s\n, got %s", expectVaule, indicator.Templ)
	}
	name = "tomcat_thread_count"
	indicator = RuleIndicator(name)
	expectVaule = `%s(tomcat_threadpool_currentthreadcount{container_label_DM_APP_ID="%s", container_label_DM_LOG_TAG!="ignore", id=~"/docker/.*", name=~"mesos.*"}) by( container_label_DM_APP_ID, container_label_DM_SLOT_INDEX ) keep_common %s %s`
	if indicator.Templ != expectVaule {
		t.Errorf("expect the templ is %s\n, got %s", expectVaule, indicator.Templ)
	}
}
