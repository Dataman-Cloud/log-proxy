package datastore

import (
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateOrIncreaseEvent(event *models.Event) error {
	var result models.Event
	notfound := db.Where("ack = ? AND severity = ? AND v_cluster = ? AND app = ? AND slot = ? AND container_id = ? AND alert_name = ?",
		false,
		event.Severity,
		event.VCluster,
		event.App,
		event.Slot,
		event.ContainerID,
		event.AlertName,
	).First(&result).RecordNotFound()
	if notfound {
		event.Count = 1
		return db.Create(event).Error
	}

	result.Count++
	result.Description = event.Description
	result.Summary = event.Summary
	return db.Save(&result).Error

}

func (db *datastore) AckEvent(ID int, userName string, groupName string) error {
	var result models.Event
	if db.Where("id = ?", ID).Where("user_name = ? OR group_name = ?", userName, groupName).First(&result).RecordNotFound() {
		return fmt.Errorf("Alert Event id=%d, user_name=%s or group_name=%s not found", ID, userName, groupName)
	}
	result.Ack = true
	return db.Save(&result).Error
}

func (db *datastore) ListAckedEvent(page models.Page, userName string, groupName string) map[string]interface{} {
	var (
		result []*models.Event
		count  int
	)
	db.Where("ack = ?", true).Where("user_name = ? OR group_name = ?", userName, groupName).Find(&result).Count(&count)
	db.Where("ack = ?", true).Where("user_name = ? OR group_name = ?", userName, groupName).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)

	return map[string]interface{}{"count": count, "events": result}
}

func (db *datastore) ListUnackedEvent(page models.Page, username string, groupname string) map[string]interface{} {
	var (
		result []*models.Event
		count  int
	)
	db.Debug().Where("ack = ?", false).Where("user_name = ? OR group_name = ?", username, groupname).Find(&result).Count(&count)
	db.Where("ack = ?", false).Where("user_name = ? OR group_name = ?", username, groupname).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	return map[string]interface{}{"count": count, "events": result}
}
