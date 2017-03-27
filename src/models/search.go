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

type LogAlertEvent struct {
	ID          uint64    `json:"id" gorm:"not null; auto_increment"`
	ContainerID string    `json:"containerid"`
	Message     string    `json:"message" binding:"required" gorm:"type:longtext"`
	LogTime     time.Time `json:"logtime" binding:"required"`
	Path        string    `json:"path" binding:"required"`
	Offset      int64     `json:"offset" binding:"required"`
	App         string    `json:"appid" binding:"required"`
	User        string    `json:"user"`
	Task        string    `json:"taskid" binding:"required"`
	Group       string    `json:"group"`
	Cluster     string    `json:"clusterid" binding:"required"`
	Keyword     string    `json:"keyword"`
}

type LogAlertRule struct {
	ID        uint64    `json:"id" gorm:"not null; auto_increment"`
	Cluster   string    `json:"cluster"`
	App       string    `json:"app"`
	Keyword   string    `json:"keyword"`
	Source    string    `json:"source"`
	Status    string    `json:"status" sql:"DEFAULT:'Enabled'"`
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
