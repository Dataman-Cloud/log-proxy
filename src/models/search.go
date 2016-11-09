package models

type Log struct {
	Message   string
	Host      string
	Port      uint64
	AppId     string
	ClusterId string
	GroupId   uint64
	Id        string
	Offset    uint64
	Path      string
	TaskId    string
}
