package alertservice

import (
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

// ReceiveAlertEvent recive the alerts from Alertmanager
func (alert *Alert) ReceiveAlertEvent(message map[string]interface{}) error {
	var err error
	for _, item := range message["alerts"].([]interface{}) {
		labels := item.(map[string]interface{})["labels"].(map[string]interface{})
		annotations := item.(map[string]interface{})["annotations"].(map[string]interface{})
		event := &models.Event{
			AlertName:     labels["alertname"].(string),
			Group:         labels["group"].(string),
			App:           labels["app"].(string),
			Indicator:     labels["indicator"].(string),
			Severity:      labels["severity"].(string),
			Task:          labels["container_label_DM_SLOT_INDEX"].(string),
			Judgement:     labels["judgement"].(string),
			ContainerID:   labels["id"].(string),
			ContainerName: labels["name"].(string),
			Value:         labels["value"].(string),
			Description:   annotations["description"].(string),
			Summary:       annotations["summary"].(string),
		}
		if err = alert.Store.CreateOrIncreaseEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (alert *Alert) GetAlertEvents(page models.Page, options map[string]interface{}) (map[string]interface{}, error) {
	result, err := alert.Store.ListEvents(page, options)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AckAlertEvent mark the alert evnet ACK
func (alert *Alert) AckAlertEvent(id int, options map[string]interface{}) error {
	switch action := options["action"].(string); action {
	case "ack":
		// TODO ugly code
		var group, app string
		if options["group"] != nil {
			group = options["group"].(string)
		}
		if options["app"] != nil {
			app = options["app"].(string)
		}
		if err := alert.Store.AckEvent(id, group, app); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("Got the wrong action options")
	}
}
