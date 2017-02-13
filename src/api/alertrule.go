package api

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

type Alert struct {
	Store       store.Store
	HTTPClient  *http.Client
	PromeServer string
}

func NewAlert() *Alert {
	return &Alert{
		Store:       datastore.InitDB(config.GetConfig().DbDriver, config.GetConfig().DbDSN),
		HTTPClient:  http.DefaultClient,
		PromeServer: config.GetConfig().PrometheusURL,
	}
}

// CreateAlertRule create the alert rule in Database
func (alert *Alert) CreateAlertRule(ctx *gin.Context) {

	var rule *models.Rule
	var err error

	if err = ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	if v := rule.Name; v == "" {
		utils.ErrorResponse(ctx, errors.New("not found Name string"))
		return
	}

	if v := rule.Alert; v == "" {
		utils.ErrorResponse(ctx, errors.New("not found Alert string"))
		return
	}

	err = alert.Store.ValidataRule(rule)
	if err == nil {
		utils.ErrorResponse(ctx, errors.New("The rule is in Database"))
		return
	}
	err = alert.Store.CreateAlertRule(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	err = alert.WriteAlertFile(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, "success")
}

func (alert *Alert) DeleteAlertRule(ctx *gin.Context) {
	var rowsAffected int64
	var ruleID uint64
	var err error
	var result models.Rule

	if ctx.Param("id") != "" {
		ruleID, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
		rowsAffected, err = alert.Store.DeleteAlertRule(ruleID)
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
		if rowsAffected == 0 {
			utils.ErrorResponse(ctx, errors.New("no rule was deleted"))
		}
		result, err = alert.Store.GetAlertRule(ruleID)
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
	}

	if ctx.Param("name") != "" && ctx.Param("alert") != "" {
		rowsAffected, err = alert.Store.DeleteAlertRuleByName(ctx.Param("name"), ctx.Param("alert"))
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
		if rowsAffected == 0 {
			utils.ErrorResponse(ctx, errors.New("no rule was deleted"))
		}
		result, err = alert.Store.GetAlertRuleByName(ctx.Param("name"), ctx.Param("alert"))
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
	}

	err = alert.RemoveAlertFile(result)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, "success")
}

func (alert *Alert) ListAlertRules(ctx *gin.Context) {
	var data []*models.Rule
	var err error

	data, err = alert.Store.ListAlertRules()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

func (alert *Alert) GetAlertRule(ctx *gin.Context) {
	var data models.Rule
	var err error
	var ruleID uint64

	if ctx.Param("id") != "" {
		ruleID, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
		data, err = alert.Store.GetAlertRule(ruleID)
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
	}
	if ctx.Param("name") != "" && ctx.Param("alert") != "" {
		data, err = alert.Store.GetAlertRuleByName(ctx.Param("name"), ctx.Param("alert"))
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
	}
	utils.Ok(ctx, data)
}

// UpdateAlertRule create the alert rule in Database
func (alert *Alert) UpdateAlertRule(ctx *gin.Context) {

	var rule *models.Rule
	var err error

	if err = ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	fmt.Printf("name: %s, alert: %s\n", rule.Name, rule.Alert)
	err = alert.Store.ValidataRule(rule)
	if err != nil {
		utils.ErrorResponse(ctx, errors.New("invalidata Rule"))
		return
	}
	err = alert.Store.UpdateAlertRule(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	err = alert.WriteAlertFile(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, "success")
}

func (alert *Alert) ReloadAlertRuleConf(ctx *gin.Context) {
	var err error

	err = alert.ReloadPrometheusConf()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, "success")
}

const ruleTempl = `# This rule was updated from {{ .UpdatedAt }}
ALERT {{.Alert}}
  IF {{ .Expr }}
  FOR {{ .Duration }}
  LABELS {{ .Labels }}
  ANNOTATIONS {description="{{ .Description }}", summary="{{ .Summary }}"}
`

func (alert *Alert) WriteAlertFile(rule *models.Rule) error {
	path := config.GetConfig().RuleFilePath

	alertfile := fmt.Sprintf("%s/%s-%s.rules", path, rule.Name, rule.Alert)
	f, err := os.Create(alertfile)
	defer f.Close()
	if err != nil {
		return err
	}

	t := template.Must(template.New("ruleTempl").Parse(ruleTempl))
	var buf bytes.Buffer
	err = t.Execute(&buf, rule)
	if err != nil {
		fmt.Println("executing templta: ", err)
		return err
	}

	f.WriteString(buf.String())

	err = alert.ReloadPrometheusConf()
	if err != nil {
		return err
	}

	return nil
}

func (alert *Alert) RemoveAlertFile(rule models.Rule) error {
	path := config.GetConfig().RuleFilePath
	alertfile := fmt.Sprintf("%s/%s-%s.rules", path, rule.Name, rule.Alert)
	f, err := os.Create(alertfile)
	defer f.Close()
	if err != nil {
		return err
	}

	message := fmt.Sprintf("# removed this rule from %s", rule.DeletedAt)
	fmt.Println("message: ", message)
	f.WriteString(message + "\n")

	err = alert.ReloadPrometheusConf()
	if err != nil {
		return err
	}

	return nil
}

func (alert *Alert) ReloadPrometheusConf() error {
	ReloadPath := "/-/reload"

	u, err := url.Parse(alert.PromeServer)
	if err != nil {
		return err
	}
	u.Path = strings.TrimRight(u.Path, "/") + ReloadPath
	resp, err := alert.HTTPClient.Post(u.String(), "application/json", nil)
	if err != nil || resp.StatusCode != 200 {
		err = fmt.Errorf("Failed to reload the configuration file of prometheus %s, return %d", u.String(), resp.StatusCode)
		return err
	}
	defer resp.Body.Close()

	return nil
}
