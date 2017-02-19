package models

import "github.com/jinzhu/gorm"

type Event struct {
	gorm.Model
	Count       int    `json:"count"`
	Severity    string `json:"severity"`
	VCluster    string `json:"vcluster"`
	App         string `json:"app"`
	Slot        string `json:"slot"`
	ContainerId string `json:"container_id"`
	AlertName   string `json:"alert_name"`
	Ack         bool   `json:"ack";sql:"DEFAULT:'false'"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
