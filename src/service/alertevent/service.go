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

func (manager *EventManager) CreateOrIncreaseEvent(event *models.Event) error {
	var result models.Event
	notfound := manager.DB.Where("ack = ? AND severity = ? AND v_cluster = ? AND app = ? AND slot = ? AND container_id = ? AND alert_name = ?",
		false,
		event.Severity,
		event.VCluster,
		event.App,
		event.Slot,
		event.ContainerId,
		event.AlertName,
	).First(&result).RecordNotFound()
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
	var result models.Event
	if manager.DB.First(&result, pk).RecordNotFound() {
		return errors.New(fmt.Sprintf("Alert Event id=%d not found", pk))
	}
	result.Ack = true
	return manager.DB.Save(&result).Error
}

func (manager *EventManager) ListAckedEvent(page models.Page) map[string]interface{} {
	var result []*models.Event
	var count int
	manager.DB.Where("ack = ?", true).Find(&result).Count(&count)
	manager.DB.Where("ack = ?", true).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	return map[string]interface{}{"count": count, "events": result}
}

func (manager *EventManager) ListUnackedEvent(page models.Page) map[string]interface{} {
	var result []*models.Event
	var count int
	manager.DB.Where("ack = ?", false).Find(&result).Count(&count)
	manager.DB.Where("ack = ?", false).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	return map[string]interface{}{"count": count, "events": result}
}
