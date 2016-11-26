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
	Period     int64  `json:"period"`
	AppId      string `json:"appid"`
	Keyword    string `json:"keyword"`
	Condition  int64  `json:"condition"`
	Enable     bool   `json:"enable"`
	CreateTime string `json:"createtime"`
}

type KeywordAlertHistory struct {
	Id         string `json:"id,omitempty"`
	AppId      string `json:"appid"`
	Keyword    string `json:"keyword"`
	Count      int64  `json:"count"`
	CreateTime string `json:"createtime"`
}
