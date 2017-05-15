package prometheusrule

import (
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

var appLabel = config.LogAppLabel()
var slotLabel = config.LogSlotLabel()

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

const lablePrefix = "container_label_"

func RuleIndicator(name string) *models.Indicator {
	var filter, appLabel, slotLabel, unit, templ string
	appLabel = fmt.Sprintf("%s%s", lablePrefix, config.LogAppLabel())
	slotLabel = fmt.Sprintf("%s%s", lablePrefix, config.LogSlotLabel())

	baseFilter := "id=~\"/docker/.*\", name=~\"mesos.*\""
	filter = appLabel + "=\"%s\", container_label_DM_LOG_TAG!=\"ignore\", " + baseFilter
	byLabels := fmt.Sprintf("( %s, %s )", appLabel, slotLabel)
	switch name {
	case "cpu_usage":
		templ = "%s(irate(container_cpu_usage_seconds_total{" + filter + "}[%s])) by" + byLabels + " keep_common * 100 %s %s"
		unit = "%"
	case "mem_usage":
		templ = "%s(container_memory_usage_bytes{" + filter +
			"} / container_spec_memory_limit_bytes{" + filter + "}) by" + byLabels + " keep_common * 100 %s %s"
		unit = "%"
	case "tomcat_thread_count":
		templ = "%s(tomcat_threadpool_currentthreadcount{" + filter + "}) by" + byLabels + " keep_common %s %s"
		unit = ""
	}

	indicator := &models.Indicator{
		Name:  name,
		Alias: name,
		Unit:  unit,
		Templ: templ,
	}
	return indicator
}
