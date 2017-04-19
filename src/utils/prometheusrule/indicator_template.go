package prometheusrule

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

var filter = "id=~\"/docker/.*\", name=~\"mesos.*\", container_label_VCLUSTER=\"%s\", container_label_APP_ID=\"%s\", container_label_DM_LOG_TAG!=\"ignore\""

var cpuUsagePercent = &models.Indicator{
	Name:  "cpu_usage",
	Templ: "%s(irate(container_cpu_usage_seconds_total{" + filter + "}[%s])) by (container_label_APP_ID) keep_common * 100 %s %s",
	Unit:  "%",
}

var memUsagePercent = &models.Indicator{
	Name: "mem_usage",
	Templ: "%s(container_memory_usage_bytes{" + filter +
		"} / container_spec_memory_limit_bytes{" + filter + "}) by (container_label_APP_ID) keep_common * 100 %s %s",
	Unit: "%",
}

var tomcatThreadPool = &models.Indicator{
	Name:  "tomcat_thread_count",
	Templ: "%s(tomcat_threadpool_currentthreadcount{" + filter + "}) by (container_label_APP_ID) keep_common %s %s",
	Unit:  "个",
}

func GetRuleIndicators() map[string]*models.Indicator {
	var ruleIndicators = make(map[string]*models.Indicator)
	ruleIndicators[cpuUsagePercent.Name] = cpuUsagePercent
	ruleIndicators[memUsagePercent.Name] = memUsagePercent
	ruleIndicators[tomcatThreadPool.Name] = tomcatThreadPool
	return ruleIndicators
}
