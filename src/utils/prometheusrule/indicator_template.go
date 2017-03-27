package prometheusrule

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

var filter = "id=~\"/docker/.*\", name=~\"mesos.*\", container_label_VCLUSTER=\"%s\", container_label_APP_ID=\"%s\""

var cpuUsagePercent = &models.Indicator{
	Name:  "cpu_usage",
	Templ: "%s(rate(container_cpu_usage_seconds_total{" + filter + "}[%s])) by (container_env_mesos_task_id) keep_common * 100 %s %s",
	Unit:  "%",
}

var memUsagePercent = &models.Indicator{
	Name: "mem_usage",
	Templ: "%s(container_memory_usage_bytes{" + filter +
		"} / container_spec_memory_limit_bytes{" + filter + "}) by (container_env_mesos_task_id) keep_common * 100 %s %s",
	Unit: "%",
}

func GetRuleIndicators() map[string]*models.Indicator {
	var ruleIndicators = make(map[string]*models.Indicator)
	ruleIndicators[cpuUsagePercent.Name] = cpuUsagePercent
	ruleIndicators[memUsagePercent.Name] = memUsagePercent
	return ruleIndicators
}
