package prometheusrule

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

var filter = "id=~\"/docker/.*\", name=~\"mesos.*\", container_label_DM_APP_ID=\"%s\", container_label_DM_LOG_TAG!=\"ignore\""

var cpuUsagePercent = &models.Indicator{
	Name:  "cpu_usage",
	Alias: "cpu_usage",
	Templ: "%s(irate(container_cpu_usage_seconds_total{" + filter + "}[%s])) by (container_label_DM_APP_ID, container_label_DM_SLOT_INDEX) keep_common * 100 %s %s",
	Unit:  "%",
}

var memUsagePercent = &models.Indicator{
	Name:  "mem_usage",
	Alias: "mem_usage",
	Templ: "%s(container_memory_usage_bytes{" + filter +
		"} / container_spec_memory_limit_bytes{" + filter + "}) by (container_label_DM_APP_ID, container_label_DM_SLOT_INDEX) keep_common * 100 %s %s",
	Unit: "%",
}

var tomcatThreadPool = &models.Indicator{
	Name:  "tomcat_thread_count",
	Alias: "tomcat_thread_count",
	Templ: "%s(tomcat_threadpool_currentthreadcount{" + filter + "}) by (container_label_DM_APP_ID, container_label_DM_SLOT_INDEX) keep_common %s %s",
	Unit:  "",
}
