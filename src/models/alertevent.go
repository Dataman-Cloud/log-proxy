package models

import "github.com/jinzhu/gorm"

type Event struct {
	gorm.Model
	Count         int    `json:"count"`
	Severity      string `json:"severity"`
	Indicator     string `json:"indicator"`
	Cluster       string `json:"cluster"`
	App           string `json:"app"`
	Task          string `json:"task"`
	Judgement     string `json:"judgement"`
	UserName      string `json:"user_name"`
	GroupName     string `json:"group_name"`
	ContainerID   string `json:"container_id"`
	ContainerName string `json:"container_name"`
	AlertName     string `json:"alert_name"`
	Ack           bool   `json:"ack";sql:"DEFAULT:'false'"`
	Value         string `json:"value"`
	Description   string `json:"description"`
	Summary       string `json:"summary"`
}
