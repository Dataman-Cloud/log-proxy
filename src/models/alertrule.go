package models

import "time"

type RuleOperation struct {
	Rule *RawRule
	File string
	MD5  []byte
}

func NewRuleOperation() *RuleOperation {
	return &RuleOperation{}
}

type RawRule struct {
	UpdatedAt   time.Time `json:"UpdatedAt"`
	Alert       string    `json:"alert"`
	Expr        string    `json:"if"`
	Pending     string    `json:"for"`
	Labels      string    `json:"labels"`
	Annotations string    `json:"annotations"`
}

type Rule struct {
	ID          uint64    `json:"ID" gorm:"primary_key"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	Class       string    `json:"class" gorm:"not null;unique_index:idx_rule"`
	Name        string    `json:"name" gorm:"not null;unique_index:idx_rule"`
	Status      string    `json:"status"`
	Cluster     string    `json:"cluster" gorm:"not null;unique_index:idx_rule"`
	App         string    `json:"app" gorm:"not null;unique_index:idx_rule"`
	User        string    `json:"user"`
	UserGroup   string    `json:"user_group"`
	Pending     string    `json:"pending"`
	Duration    string    `json:"duration"`
	Indicator   string    `json:"indicator"`
	Severity    string    `json:"severity"`
	Aggregation string    `json:"aggregation"` // max, min, avg, sum, count
	Comparison  string    `json:"comparison"`
	Threshold   int64     `json:"threshold"`
}

func NewRule() *Rule {
	return &Rule{
		Class:     "mola",
		Status:    "Uninit", //uninit, Enabled, Disabled
		Duration:  "5m",
		Threshold: int64(0),
	}
}

type Indicator struct {
	Name  string `json:"name"`
	Templ string `json:"template"`
	Unit  string `json:"unit"`
}
