package models

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
