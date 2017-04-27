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
	Cluster     string    `json:"cluster"`
	App         string    `json:"app"`
	Keyword     string    `json:"keyword"`
	Source      string    `json:"source"`
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
