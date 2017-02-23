package datastore

import (
	"errors"
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateOrIncreaseEvent(event *models.Event) error {
	var result models.Event
	notfound := db.DB.Where("ack = ? AND severity = ? AND v_cluster = ? AND app = ? AND slot = ? AND container_id = ? AND alert_name = ?",
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
		return db.DB.Create(event).Error
	} else {
		result.Count += 1
		result.Description = event.Description
		result.Summary = event.Summary
		return db.DB.Save(&result).Error
	}
}

func (db *datastore) AckEvent(pk int, username string, groupname string) error {
	var result models.Event
	if db.DB.Where("id = ?", pk).Where("user_name = ? OR group_name = ?", username, groupname).First(&result).RecordNotFound() {
		return errors.New(fmt.Sprintf("Alert Event id=%d, user_name=%s or group_name=%s not found", pk, username, groupname))
	}
	result.Ack = true
	return db.DB.Save(&result).Error
}

func (db *datastore) ListAckedEvent(page models.Page, username string, groupname string) map[string]interface{} {
	var (
		result []*models.Event
		count  int
	)
	db.DB.Where("ack = ?", true).Where("user_name = ? OR group_name = ?", username, groupname).Find(&result).Count(&count)
	db.DB.Where("ack = ?", true).Where("user_name = ? OR group_name = ?", username, groupname).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)

	return map[string]interface{}{"count": count, "events": result}
}

func (db *datastore) ListUnackedEvent(page models.Page, username string, groupname string) map[string]interface{} {
	var (
		result []*models.Event
		count  int
	)
	db.DB.Where("ack = ?", false).Where("user_name = ? OR group_name = ?", username, groupname).Find(&result).Count(&count)
	db.DB.Where("ack = ?", false).Where("user_name = ? OR group_name = ?", username, groupname).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	return map[string]interface{}{"count": count, "events": result}
}
