package models

import "github.com/jinzhu/gorm"

type AlertEvent struct {
	gorm.Model
	Count       int    `json:"count"`
	Severity    string `json:"severity";gorm:"unique_index:union"`
	VCluster    string `json:"vcluster";gorm:"unique_index:union"`
	App         string `json:"app";gorm:"unique_index:union"`
	Slot        string `json:"slot";gorm:"unique_index:union"`
	ContainerId string `json:"container_id";gorm:"unique_index:union`
	AlertName   string `json:"alert_name";gorm:"unique_index:union"`
	Ack         bool   `json:"ack";sql:"DEFAULT:'false'"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
