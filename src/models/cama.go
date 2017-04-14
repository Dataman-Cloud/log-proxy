package models

type CamaEvent struct {
	ID        string `json:"IDENTIFIER"`
	Channel   string `json:"CHANNEL"`
	FirstTime string `json:"FIRST_TIME"`
	LastTime  string `json:"LAST_TIME"`
	Recover   int    `json:"RECOVER"`
	Merger    int    `json:"MERGER"`
	Node      string `json:"NODE"`
	NodeAlias string `json:"NODEALIAS"`
	ServerNo  string `json:"SERVER_NO"`
	EventDesc string `json:"EVENT_DESC"`
	Level     int    `json:"LEVEL"`
}
