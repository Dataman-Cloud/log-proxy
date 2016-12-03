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

type Alert struct {
	Id         string `json:"id,omitempty"`
	AppId      string `json:"appid"`
	Keyword    string `json:"keyword"`
	Path       string `json:"path"`
	CreateTime string `json:"createtime"`
}
