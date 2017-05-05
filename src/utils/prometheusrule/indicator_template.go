package prometheusrule

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

var filter = "id=~\"/docker/.*\", name=~\"mesos.*\", container_label_DM_APP_ID=\"%s\", container_label_DM_LOG_TAG!=\"ignore\""

var cpuUsagePercent = &models.Indicator{
	Name:  "cpu_usage",
	Alias: "CPU使用百分比",
	Templ: "%s(irate(container_cpu_usage_seconds_total{" + filter + "}[%s])) by (container_label_DM_APP_ID) keep_common * 100 %s %s",
	Unit:  "%",
}

var memUsagePercent = &models.Indicator{
	Name:  "mem_usage",
	Alias: "内存使用百分比",
	Templ: "%s(container_memory_usage_bytes{" + filter +
		"} / container_spec_memory_limit_bytes{" + filter + "}) by (container_label_DM_APP_ID) keep_common * 100 %s %s",
	Unit: "%",
}

var tomcatThreadPool = &models.Indicator{
	Name:  "tomcat_thread_count",
	Alias: "Tomcat线程数",
	Templ: "%s(tomcat_threadpool_currentthreadcount{" + filter + "}) by (container_label_DM_APP_ID) keep_common %s %s",
	Unit:  "",
}
