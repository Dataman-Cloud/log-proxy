package camaalert

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/prometheus/common/log"
)

const camaEventTempl = `集群{{ .Cluster }}应用{{ .App }}的{{ .Indicator }}{{ .Judgement }}{{.Operator}}{{.JudgeValue}}`

func Event2Cama(event *models.Event) *models.CamaEvent {
	var recover int
	if event.Ack == true {
		recover = 1
	} else {
		recover = 0
	}
	camaEvent := &models.CamaEvent{
		ID:        event.AlertName,
		Channel:   "DOCKER",
		FirstTime: event.CreatedAt,
		LastTime:  event.UpdatedAt,
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
		Cluster:   event.Cluster,
		App:       event.App,
		Indicator: event.Indicator,
	}
	camaEventDesc.Judgement = strings.Split(event.Judgement, " ")[0]
	camaEventDesc.Operator = strings.Split(event.Judgement, " ")[1]
	camaEventDesc.JudgeValue = strings.Split(event.Judgement, " ")[2]

	t := template.Must(template.New("camaEventTempl").Parse(camaEventTempl))
	var buf bytes.Buffer
	err = t.Execute(&buf, camaEventDesc)
	if err != nil {
		log.Errorln("executing templta: ", err)
		return err
	}

	return buf.String()
}
