package alertevent

import (
	"errors"
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/jinzhu/gorm"
)

type EventManager struct {
	DB *gorm.DB
}

func (manager *EventManager) CreateOrIncreaseEvent(event *models.AlertEvent) error {
	var result models.AlertEvent
	notfound := manager.DB.Where("ack = ? AND severity = ? AND v_cluster = ? AND app = ? AND slot = ? AND container_id = ? AND alert_name = ?",
		false,
		event.Severity,
		event.VCluster,
		event.App,
		event.Slot,
		event.ContainerId,
		event.AlertName,
	).First(&result).RecordNotFound()
	/* TODO Query With Where (Struct & Map) work out, why Ack(boolean) ignored?
	notfound := manager.DB.Debug().Where(&models.AlertEvent{
		Ack:         false,
		Severity:    event.Severity,
		VCluster:    event.VCluster,
		App:         event.App,
		Slot:        event.Slot,
		ContainerId: event.ContainerId,
		AlertName:   event.AlertName,
	}).First(&result).RecordNotFound()
	*/
	if notfound {
		event.Count = 1
		return manager.DB.Create(event).Error
	} else {
		result.Count += 1
		result.Description = event.Description
		result.Summary = event.Summary
		return manager.DB.Save(&result).Error
	}
}

func (manager *EventManager) AckEvent(pk int) error {
	var result models.AlertEvent
	if manager.DB.First(&result, pk).RecordNotFound() {
		return errors.New(fmt.Sprintf("Alert Event id=%d not found", pk))
	}
	result.Ack = true
	return manager.DB.Save(&result).Error
}

func (manager *EventManager) ListAckedEvent(page models.Page) []*models.AlertEvent {
	var result []*models.AlertEvent
	manager.DB.Where("ack = ?", true).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	return result
}

func (manager *EventManager) ListUnackedEvent(page models.Page) []*models.AlertEvent {
	var result []*models.AlertEvent
	manager.DB.Where("ack = ?", false).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	return result
}
