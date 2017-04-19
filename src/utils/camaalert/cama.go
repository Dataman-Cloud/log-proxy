package camaalert

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"

	log "github.com/Sirupsen/logrus"
)

const camaEventTempl = `集群 {{ .Cluster }} 应用 {{ .App }} 的 {{ .Indicator }} {{ .Judgement }} {{.Operator}} {{.JudgeValue}}`

func Event2CamaEvent(event *models.Event) *models.CamaEvent {
	var recover int
	if event.Ack == true {
		recover = 1
	} else {
		recover = 0
	}
	camaEvent := &models.CamaEvent{
		ID:        event.AlertName,
		Channel:   "DOCKER",
		FirstTime: event.CreatedAt.Format(config.CamaTimeFormatString),
		LastTime:  event.UpdatedAt.Format(config.CamaTimeFormatString),
		Recover:   recover,
		Merger:    event.Count,
		Node:      "",
		NodeAlias: "",
		ServerNo:  "",
		EventDesc: Event2Desc(event),
		Level:     5,
	}
	return camaEvent
}

func Event2Desc(event *models.Event) string {
	camaEventDesc := &models.CamaEventDesc{
		Cluster: event.Cluster,
		App:     event.App,
	}
	indicator := event.Indicator
	switch indicator {
	case "cpu_usage":
		camaEventDesc.Indicator = "CPU使用率"
	case "mem_usage":
		camaEventDesc.Indicator = "内存占用率"
	case "tomcat_thread_count":
		camaEventDesc.Indicator = "Tomcat线程数"
	}

	judgement := strings.Split(event.Judgement, " ")
	if len(judgement) != 3 {
		log.Errorf("Failed to convert the event to cama, Judgement is %s", event.Judgement)
		camaEventDesc.Judgement = "jugement"
		camaEventDesc.Operator = "operator"
		camaEventDesc.JudgeValue = "value"
	} else {

		judge := strings.Split(event.Judgement, " ")[0]
		switch judge {
		case "max":
			camaEventDesc.Judgement = "最大值"
		case "min":
			camaEventDesc.Judgement = "最小值"
		case "avg":
			camaEventDesc.Judgement = "平均值"
		case "sum":
			camaEventDesc.Judgement = "总和"
		}
		operater := strings.Split(event.Judgement, " ")[1]
		switch operater {
		case ">":
			camaEventDesc.Operator = "大于"
		case "<":
			camaEventDesc.Operator = "小于"
		case "==":
			camaEventDesc.Operator = "等于"
		}
		camaEventDesc.JudgeValue = strings.Split(event.Judgement, " ")[2]
	}
	t := template.Must(template.New("camaEventTempl").Parse(camaEventTempl))
	var buf bytes.Buffer
	err := t.Execute(&buf, camaEventDesc)
	if err != nil {
		log.Errorln("executing templta: ", err)
	}

	return buf.String()
}
