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
