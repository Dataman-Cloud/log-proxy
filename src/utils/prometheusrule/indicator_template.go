package prometheusrule

import (
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func RuleIndicator(name string) *models.Indicator {
	var filter, appLabel, slotLabel, tagLabel, unit, templ string
	appLabel = config.MonitorAppLabel()
	slotLabel = config.MonitorSlotLabel()
	tagLabel = config.LogTagLabel()

	baseFilter := "id=~\"/docker/.*\", name=~\"mesos.*\""
	filter = appLabel + "=\"%s\", " + tagLabel + "!=\"ignore\", " + baseFilter
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
