package models

import "github.com/jinzhu/gorm"

type Event struct {
	gorm.Model
	AlertName     string `json:"alert_name" gorm:"not null;column:alertname"`
	Group         string `json:"group" gorm:"not null;column:groupname"`
	App           string `json:"app"`
	Task          string `json:"task"`
	Severity      string `json:"severity"`
	Indicator     string `json:"indicator"`
	Judgement     string `json:"judgement"`
	ContainerID   string `json:"container_id"`
	ContainerName string `json:"container_name"`
	Ack           bool   `json:"ack"`
	Value         string `json:"value"`
	Description   string `json:"description"`
	Summary       string `json:"summary"`
	Count         int    `json:"count"`
}
