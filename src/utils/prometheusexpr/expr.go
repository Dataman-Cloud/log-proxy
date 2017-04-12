package prometheusexpr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Dataman-Cloud/log-proxy/src/config"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var exprs map[string]*Expr

// Expr define the prometheus expr
type Expr struct {
	Name        string  `json:"name"`
	Aggregation string  `json:"aggregation"`
	Function    string  `json:"function"`
	Metric      string  `json:"metric"`
	Filter      *Filter `json:"filter"`
	By          *By     `json:"by"`
}

func NewExpr() *Expr {
	filter := &Filter{}
	by := &By{}
	return &Expr{
		Filter: filter,
		By:     by,
	}
}

// Filter define the fileter labels
type Filter struct {
	Cluster string `json:"cluster"`
	User    string `json:"user"`
	App     string `json:"app"`
	Slot    string `json:"slot"`
	Task    string `json:"task"`
	Fixed   string `json:"fixed"`
}

// GetFilter return the default filter strings
func GetFilter() *Filter {
	var prefix string
	if config.GetConfig() != nil {
		prefix = config.GetConfig().QueryPrefix
	} else {
		prefix = "DM"
	}
	cluster := fmt.Sprintf("container_label_%s_VCLUSTER", prefix)
	user := fmt.Sprintf("container_label_%s_USER", prefix)
	app := fmt.Sprintf("container_label_%s_APP_NAME", prefix)
	slot := fmt.Sprintf("container_label_%s_SLOT_ID", prefix)
	task := fmt.Sprintf("container_label_%s_TASK_ID", prefix)
	filter := &Filter{
		Cluster: cluster + "='%s'",
		User:    user + "='%s'",
		App:     app + "='%s'",
		Slot:    slot + "='%s'",
		Task:    task + "='%s'",
		Fixed:   "id=~'/docker/.*', name=~'mesos.*', container_label_DM_LOG_TAG!='ignore'",
	}
	return filter
}

// By define the by labels
type By struct {
	Cluster string `json:"cluster"`
	User    string `json:"user"`
	App     string `json:"app"`
	Slot    string `json:"slot"`
	Task    string `json:"task"`
}

// GetBy return the default By strings
func GetBy() *By {
	var prefix string
	if config.GetConfig() != nil {
		prefix = config.GetConfig().QueryPrefix
	} else {
		prefix = "DM"
	}
	cluster := fmt.Sprintf("container_label_%s_VCLUSTER", prefix)
	user := fmt.Sprintf("container_label_%s_USER", prefix)
	app := fmt.Sprintf("container_label_%s_APP_NAME", prefix)
	slot := fmt.Sprintf("container_label_%s_SLOT_ID", prefix)
	task := fmt.Sprintf("container_label_%s_TASK_ID", prefix)
	by := &By{
		Cluster: cluster,
		User:    user,
		App:     app,
		Slot:    slot,
		Task:    task,
	}
	return by
}

// Exprs get the exprs from files
func Exprs(path string) error {
	if len(exprs) == 0 {
		exprs = make(map[string]*Expr, 0)
		files, err := getExprFiles(path)
		if err != nil {
			return err
		}
		log.Infof("Got the %d expr files", len(files))

		for _, file := range files {
			expr, err := parseExpr(path + "/" + file)
			if err != nil {
				return err
			}
			exprs[expr.Name] = expr
		}
	}
	return nil
}

// GetExprs init the Query
func GetExprs() map[string]*Expr {
	return exprs
}

func parseExpr(file string) (*Expr, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	expr := NewExpr()
	err = yaml.Unmarshal(fileContent, &expr)
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, fmt.Errorf("parse the file %s error", file)
	}

	if config.GetConfig() != nil {
		expr.setFilterAndByWithPrefix(config.GetConfig().QueryPrefix)
	} else {
		expr.setFilterAndByWithPrefix("DM")
	}

	return expr, nil
}

func getExprFiles(path string) ([]string, error) {
	log.Infof("Getting the expr files in dir %s", path)
	files := make([]string, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		files = append(files, f.Name())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, err
}

func (expr *Expr) setFilterAndByWithPrefix(prefix string) {
	filter := GetFilter()

	if expr.Filter.Cluster == "" {
		expr.Filter.Cluster = filter.Cluster
	}
	if expr.Filter.User == "" {
		expr.Filter.User = filter.User
	}
	if expr.Filter.App == "" {
		expr.Filter.App = filter.App
	}
	if expr.Filter.Slot == "" {
		expr.Filter.Slot = filter.Slot
	}
	if expr.Filter.Task == "" {
		expr.Filter.Task = filter.Task
	}
	if expr.Filter.Fixed == "" {
		expr.Filter.Fixed = filter.Fixed
	}

	by := GetBy()
	if expr.By.Cluster == "" {
		expr.By.Cluster = by.Cluster
	}
	if expr.By.User == "" {
		expr.By.User = by.User
	}
	if expr.By.App == "" {
		expr.By.App = by.App
	}
	if expr.By.Slot == "" {
		expr.By.Slot = by.Slot
	}
	if expr.By.Task == "" {
		expr.By.Task = by.Task
	}
}
