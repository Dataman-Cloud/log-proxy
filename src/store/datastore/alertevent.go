package datastore

import (
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateOrIncreaseEvent(event *models.Event) error {
	var result models.Event
	notfound := db.Where("ack = ? AND severity = ? AND cluster = ? AND app = ? AND task = ? AND container_id = ? AND container_name =? AND alert_name = ?",
		false,
		event.Severity,
		event.Cluster,
		event.App,
		event.Task,
		event.ContainerID,
		event.ContainerName,
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

func (db *datastore) AckEvent(ID int, cluster, app string) error {
	var result models.Event
	if db.Where("id = ?", ID).Where("cluster = ? AND app = ?", cluster, app).First(&result).RecordNotFound() {
		return fmt.Errorf("Alert Event id=%d, cluster=%s or app=%s not found", ID, cluster, app)
	}
	result.Ack = true
	return db.Save(&result).Error
}

func (db *datastore) ListAckedEvent(page models.Page, cluster, app string) map[string]interface{} {
	var (
		result []*models.Event
		count  int
	)
	if cluster == "" && app == "" {
		db.Where("ack = ?", true).Find(&result).Count(&count)
		db.Where("ack = ?", true).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	} else if cluster != "" && app == "" {
		db.Where("ack = ?", true).Where("cluster = ?", cluster).Find(&result).Count(&count)
		db.Where("ack = ?", true).Where("cluster = ?", cluster).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	} else if cluster != "" && app != "" {
		db.Where("ack = ?", true).Where("cluster = ? AND app = ?", cluster, app).Find(&result).Count(&count)
		db.Where("ack = ?", true).Where("cluster = ? AND app = ?", cluster, app).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	} else {
		result = nil
	}
	return map[string]interface{}{"count": count, "events": result}
}

func (db *datastore) ListUnackedEvent(page models.Page, cluster, app string) map[string]interface{} {
	var (
		result []*models.Event
		count  int
	)
	if cluster == "" && app == "" {
		db.Where("ack = ?", false).Find(&result).Count(&count)
		db.Where("ack = ?", false).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	} else if cluster != "" && app == "" {
		db.Where("ack = ?", false).Where("cluster = ?", cluster).Find(&result).Count(&count)
		db.Where("ack = ?", false).Where("cluster = ?", cluster).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	} else if cluster != "" && app != "" {
		db.Where("ack = ?", false).Where("cluster = ? AND app = ?", cluster, app).Find(&result).Count(&count)
		db.Where("ack = ?", false).Where("cluster = ? AND app = ?", cluster, app).Offset(page.PageFrom).Limit(page.PageSize).Find(&result)
	} else {
		result = nil
	}
	return map[string]interface{}{"count": count, "events": result}
}
