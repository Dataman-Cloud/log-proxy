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

type LogAlertRule struct {
	ID        uint64    `json:"id" gorm:"NOT NULL; AUTO_INCREMENT"`
	App       string    `json:"app"`
	Keyword   string    `json:"keyword"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TaskInfo struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	LogCount int64  `json:"logCount"`
}

const (
	TaskRunning string = "running"
	TaskDied    string = "died"
)
