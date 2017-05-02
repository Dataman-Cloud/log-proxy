package models

import "time"

// Log log struct
type Log struct {
	Message   string
	Host      string
	Port      uint64
	AppID     string
	ClusterID string
	GroupID   uint64
	ID        string
	Offset    uint64
	Path      string
	TaskID    string
}

// Alert keyword filter struct
type Alert struct {
	ID         string `json:"id,omitempty"`
	AppID      string `json:"app"`
	Keyword    string `json:"keyword"`
	Path       string `json:"path"`
	CreateTime string `json:"createtime"`
}

type LogAlertRule struct {
	ID          uint64    `json:"id" gorm:"not null; auto_increment"`
	User        string    `json:"user" gorm:"not null" binding:"required"`
	Group       string    `json:"group" gorm:"not null" binding:"required"`
	Cluster     string    `json:"cluster" gorm:"not null" binding:"required"`
	App         string    `json:"app" gorm:"not null" binding:"required"`
	Keyword     string    `json:"keyword" gorm:"not null" binding:"required"`
	Source      string    `json:"source" gorm:"not null" binding:"required"`
	Status      string    `json:"status" sql:"DEFAULT:'Enabled'"`
	Description string    `json:"description"  gorm:"type:longtext"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const (
	TaskRunning string = "running"
	TaskDied    string = "died"
)

type LogAlertClusters struct {
	Cluster string `json:"cluster"`
}

type LogAlertApps struct {
	App string `json:"app"`
}
