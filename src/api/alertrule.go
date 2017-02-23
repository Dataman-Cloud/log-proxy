package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"

	"github.com/gin-gonic/gin"
)

const (
	// Receive Alert event error
	ReceiveEventError = "503-21000"
	// Ack Alert event error
	AckEventError = "503-21001"

	ruleTempl = `# This rule was updated from {{ .UpdatedAt }} by {{ .Name }}
  ALERT {{.Alert}}
  IF {{ .Expr }}
  FOR {{ .Duration }}
  LABELS {{ .Labels }}
  ANNOTATIONS {description="{{ .Description }}", summary="{{ .Summary }}"}
`
)

type Alert struct {
	Store       store.Store
	HTTPClient  *http.Client
	PromeServer string
}

func NewAlert() *Alert {
	return &Alert{
		Store:       datastore.From(database.GetDB()),
		HTTPClient:  http.DefaultClient,
		PromeServer: config.GetConfig().PrometheusURL,
	}
}

// CreateAlertRule create the alert rule in Database
func (alert *Alert) CreateAlertRule(ctx *gin.Context) {

	var rule *models.Rule
	var data models.Rule
	var err error

	if err = ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	if v := rule.Name; v == "" {
		utils.ErrorResponse(ctx, errors.New("not found Name string"))
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

	data, err = alert.Store.GetAlertRuleByName(rule.Name, rule.Alert)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, data)
}

func (alert *Alert) DeleteAlertRule(ctx *gin.Context) {
	var rowsAffected int64
	var err error
	var rule *models.Rule
	var result models.Rule

	if err = ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	result, err = alert.Store.GetAlertRule(rule.ID)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	rowsAffected, err = alert.Store.DeleteAlertRuleByIDName(rule.ID, rule.Name)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	if rowsAffected == 0 {
		utils.ErrorResponse(ctx, errors.New("no rule was deleted"))
	}

	err = alert.RemoveAlertFile(result)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, "success")
}

func (alert *Alert) ListAlertRules(ctx *gin.Context) {
	data, err := alert.Store.ListAlertRules(ctx.MustGet("page").(models.Page), ctx.Query("name"))
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

	utils.Ok(ctx, data)
}

// UpdateAlertRule create the alert rule in Database
func (alert *Alert) UpdateAlertRule(ctx *gin.Context) {

	var rule *models.Rule
	var data models.Rule
	var err error

	if err = ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
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

	data, err = alert.Store.GetAlertRuleByName(rule.Name, rule.Alert)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, data)
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

	message := fmt.Sprintf("# removed this rule")
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

func (alert *Alert) ReceiveAlertEvent(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ReceiveEventError, err))
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ReceiveEventError, err))
		return
	}

	for _, item := range m["alerts"].([]interface{}) {
		labels := item.(map[string]interface{})["labels"].(map[string]interface{})
		annotations := item.(map[string]interface{})["annotations"].(map[string]interface{})
		event := &models.Event{
			AlertName:   labels["alertname"].(string),
			Severity:    labels["severity"].(string),
			VCluster:    labels["container_label_VCLUSTER"].(string),
			App:         labels["container_label_APP"].(string),
			Slot:        labels["container_label_SLOT"].(string),
			UserName:    labels["container_label_USER_NAME"].(string),
			GroupName:   labels["container_label_GROUP_NAME"].(string),
			ContainerID: labels["id"].(string),
			Description: annotations["description"].(string),
			Summary:     annotations["summary"].(string),
		}
		if err := alert.Store.CreateOrIncreaseEvent(event); err != nil {
			utils.ErrorResponse(ctx, utils.NewError(ReceiveEventError, err))
			return
		}
	}

	utils.Ok(ctx, map[string]string{"status": "success"})
}

func (alert *Alert) AckAlertEvent(ctx *gin.Context) {
	pk, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}
	var data map[string]interface{}
	if err := ctx.BindJSON(&data); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}

	switch action := data["action"].(string); action {
	case "ack":
		// TODO ugly code
		var username, groupname string
		if data["user_name"] != nil {
			username = data["user_name"].(string)
		}
		if data["group_name"] != nil {
			groupname = data["group_name"].(string)
		}
		if err = alert.Store.AckEvent(pk, username, groupname); err != nil {
			utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
			return
		}
		utils.Ok(ctx, map[string]string{"status": "success"})
	}
}

func (alert *Alert) GetAlertEvents(ctx *gin.Context) {
	switch ack := ctx.Query("ack"); ack {
	case "true":
		result := alert.Store.ListAckedEvent(ctx.MustGet("page").(models.Page), ctx.Query("user_name"), ctx.Query("group_name"))
		utils.Ok(ctx, result)
	case "false", "":
		result := alert.Store.ListUnackedEvent(ctx.MustGet("page").(models.Page), ctx.Query("user_name"), ctx.Query("group_name"))
		utils.Ok(ctx, result)
	}
}
