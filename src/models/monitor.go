package models

// QueryParameter define the fields of query paramter,
// Frontend pass the parameter in 2 way:
// 1. Expr, Start, End, Step
// 2. Metric, Cluster, User, App, Task, Start, End, Step
type QueryParameter struct {
	Expr    string
	Metric  string
	Cluster string
	User    string
	App     string
	Task    string
	Start   string
	End     string
	Step    string
}
