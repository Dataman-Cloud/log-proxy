package models

import "github.com/jinzhu/gorm"

type Event struct {
	gorm.Model
	Count       int    `json:"count"`
	Severity    string `json:"severity" gorm:"unique_index:idx_event"`
	VCluster    string `json:"vcluster" gorm:"unique_index:idx_event"`
	App         string `json:"app" gorm:"unique_index:idx_event"`
	Slot        string `json:"slot" gorm:"unique_index:idx_event"`
	UserName    string `json:"user_name"`
	GroupName   string `json:"group_name"`
	ContainerID string `json:"container_id" gorm:"unique_index:idx_event"`
	AlertName   string `json:"alert_name" gorm:"unique_index:idx_event"`
	Ack         bool   `json:"ack";sql:"DEFAULT:'false'"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
