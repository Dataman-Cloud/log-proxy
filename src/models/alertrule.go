package models

import "time"

type RuleOperation struct {
	Rule *Rule
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
	Name        string    `json:"name" gorm:"not null;unique"`
	TenantID    uint64    `json:"tenant" gorm:"not null;"`
	Group       string    `json:"group" gorm:"not null;"`
	App         string    `json:"app" gorm:"not null;"`
	Severity    string    `json:"severity" gorm:"not null"`
	Indicator   string    `json:"indicator" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null"`
	Pending     string    `json:"pending" gorm:"not null"`
	Duration    string    `json:"duration" gorm:"not null"`
	Aggregation string    `json:"aggregation" gorm:"not null"` // max, min, avg, sum, count
	Comparison  string    `json:"comparison" gorm:"not null"`
	Threshold   int64     `json:"threshold" gorm:"not null"`
	Unit        string    `json:"unit"`
}

func NewRule() *Rule {
	return &Rule{
		Status:    "Uninitialized", //Uninitialized, Enabled, Disabled
		Duration:  "5m",
		Threshold: int64(0),
		Group:     "Undefine",
	}
}

type Indicator struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Templ string `json:"template"`
	Unit  string `json:"unit"`
}

type RulesList struct {
	Count int64   `json:"count"`
	Rules []*Rule `json:"rules"`
}

func NewRulesList() *RulesList {
	return &RulesList{
		Rules: []*Rule{},
	}
}
