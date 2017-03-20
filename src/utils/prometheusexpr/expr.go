package prometheusexpr

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var exprs map[string]*Expr

// Expr define the prometheus expr
type Expr struct {
	Name        string   `json:"name"`
	Aggregation string   `json:"aggregation"`
	Function    string   `json:"function"`
	Metric      string   `json:"metric"`
	Filter      *Filter  `json:"filter"`
	By          []string `json:"by"`
}

func NewExpr() *Expr {
	filter := &Filter{}
	by := make([]string, 0)
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
	Task    string `json:"task"`
	Fixed   string `json:"fixed"`
}

// Exprs get the exprs from files
func Exprs(path string) map[string]*Expr {
	if len(exprs) == 0 {
		exprs = make(map[string]*Expr, 0)
		files, err := getExprFiles(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			expr, err := parseExpr(path + "/" + file)
			if err != nil {
				log.Fatal(err)
			}
			exprs[expr.Name] = expr
		}
	}
	return exprs
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
	log.Infof("Got the %d expr files", len(files))
	return files, err
}
