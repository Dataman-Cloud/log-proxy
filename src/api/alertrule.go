package api

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

const (
	// ReceiveEventError code
	ReceiveEventError = "503-21000"
	// ReceiveEventError code
	AckEventError = "503-21001"

	// PromtheusReloadPath path string
	PromtheusReloadPath = "/-/reload"

	ruleTempl = `# This rule was updated by {{ .Name }}
  ALERT {{.Alert}}
  IF {{ .Expr }}
  FOR {{ .Duration }}
  LABELS {{ .Labels }}
  ANNOTATIONS {description="{{ .Description }}", summary="{{ .Summary }}"}
`
	ruleFileUpdate = "update"
	ruleFileDelete = "delete"
	ruleInterval   = "1m"
)

type Alert struct {
	Store      store.Store
	HTTPClient *http.Client
	PromServer string
	Interval   string
	RulesPath  string
}

// NewAlert init the struct Alert
func NewAlert() *Alert {
	interval := config.GetConfig().RuleFileInterval
	if interval == "" {
		interval = ruleInterval
	}

	return &Alert{
		Store:      datastore.From(database.GetDB()),
		HTTPClient: http.DefaultClient,
		PromServer: config.GetConfig().PrometheusURL,
		Interval:   interval,
		RulesPath:  config.GetConfig().RuleFilePath,
	}
}

// CreateAlertRule create the alert rule in Database
func (alert *Alert) CreateAlertRule(ctx *gin.Context) {

	var (
		rule *models.Rule
		data models.Rule
		err  error
	)
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

	data, err = alert.Store.GetAlertRuleByName(rule.Name, rule.Alert)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	err = alert.WriteAlertFile(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, data)
}

// DeleteAlertRule delete the rule by id and name from DB and files
func (alert *Alert) DeleteAlertRule(ctx *gin.Context) {
	var (
		rowsAffected int64
		err          error
		rule         *models.Rule
		result       models.Rule
	)
	if err = ctx.BindJSON(&rule); err != nil {
		log.Errorln("DeleteAlertRule: Parse JSON rule error, ", err)
		utils.ErrorResponse(ctx, err)
		return
	}

	rule.ID, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	result, err = alert.Store.GetAlertRule(rule.ID)
	if err != nil {
		log.Errorln("DeleteAlertRule: GetAlertRule() error, ", err)
		utils.ErrorResponse(ctx, err)
		return
	}

	rowsAffected, err = alert.Store.DeleteAlertRuleByIDName(rule.ID, rule.Name)
	if err != nil {
		log.Errorln("DeleteAlertRule: DeleteAlertRuleByIDName() error, ", err)
		utils.ErrorResponse(ctx, err)
		return
	}

	if rowsAffected == 0 {
		utils.ErrorResponse(ctx, errors.New("no rule was deleted"))
	}

	err = alert.RemoveAlertFile(result)
	if err != nil {
		log.Errorln("DeleteAlertRule: Delete Alert file error, ", err)
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, "success")
}

// ListAlertRules list the rules by name with pages.
func (alert *Alert) ListAlertRules(ctx *gin.Context) {
	data, err := alert.Store.ListAlertRules(ctx.MustGet("page").(models.Page), ctx.Query("name"))
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// GetAlertRule return the info of alert rule by id
func (alert *Alert) GetAlertRule(ctx *gin.Context) {
	var (
		data   models.Rule
		err    error
		ruleID uint64
	)

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

// UpdateAlertRule update the alert rule in Database
func (alert *Alert) UpdateAlertRule(ctx *gin.Context) {
	var (
		rule *models.Rule
		data models.Rule
		err  error
	)
	if err = ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	err = alert.Store.UpdateAlertRule(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	data, err = alert.Store.GetAlertRuleByName(rule.Name, rule.Alert)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	err = alert.WriteAlertFile(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, data)
}

// ReloadAlertRuleConf tigger the conf reloading by prometheus api
func (alert *Alert) ReloadAlertRuleConf(ctx *gin.Context) {
	var err error

	err = alert.ReloadPrometheusConf()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, "success")
}

// WriteAlertFile write the alert rule to file
func (alert *Alert) WriteAlertFile(rule *models.Rule) error {
	path := alert.RulesPath
	alertfile := fmt.Sprintf("%s/%s-%s.rule", path, rule.Name, rule.Alert)
	f, err := os.Create(alertfile)
	defer f.Close()
	if err != nil {
		return err
	}

	t := template.Must(template.New("ruleTempl").Parse(ruleTempl))
	var buf bytes.Buffer
	err = t.Execute(&buf, rule)
	if err != nil {
		log.Errorln("executing templta: ", err)
		return err
	}

	f.WriteString(buf.String())

	err = alert.ReloadPrometheusConf()
	if err != nil {
		return err
	}

	return nil
}

// RemoveAlertFile remove the rule from the file.
func (alert *Alert) RemoveAlertFile(rule models.Rule) error {
	path := alert.RulesPath
	alertfile := fmt.Sprintf("%s/%s-%s.rule", path, rule.Name, rule.Alert)
	f, err := os.Create(alertfile)
	defer f.Close()
	if err != nil {
		return err
	}

	message := fmt.Sprintf("# removed this rule")
	f.WriteString(message + "\n")

	err = alert.ReloadPrometheusConf()
	if err != nil {
		return err
	}

	return nil
}

// ReloadPrometheusConf reload the conf by calling prometheus api
func (alert *Alert) ReloadPrometheusConf() error {
	u, err := url.Parse(alert.PromServer)
	if err != nil {
		return err
	}
	u.Path = strings.TrimRight(u.Path, "/") + PromtheusReloadPath
	resp, err := alert.HTTPClient.Post(u.String(), "application/json", nil)
	if err != nil || resp.StatusCode != 200 {
		err = fmt.Errorf("Failed to reload the configuration file of prometheus %s, return %d", u.String(), resp.StatusCode)
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ReceiveAlertEvent recive the alerts from Alertmanager
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

// AckAlertEvent mark the alert evnet ACK
func (alert *Alert) AckAlertEvent(ctx *gin.Context) {
	var err error
	pk, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}
	var data map[string]interface{}

	if err = ctx.BindJSON(&data); err != nil {
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

// GetAlertEvents list the alert events
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

// AlertRuleFilesMaintainer keep the rule files sync with db
func (alert *Alert) AlertRuleFilesMaintainer() {
	c := cron.New()
	interval := fmt.Sprintf("@every %s", alert.Interval)
	c.AddFunc(interval, func() { alert.UpdateAlertRuleFiles() })
	c.Start()
	log.Infof("The alert files will be check in %s", alert.Interval)
	alert.UpdateAlertRuleFiles()
}

// UpdateAlertRuleFiles update the rule files
func (alert *Alert) UpdateAlertRuleFiles() {
	var (
		rules     []*models.Rule
		files     map[string]interface{}
		ruleNames []*models.RuleOperation
		err       error
	)
	reloadReady := false

	path := alert.RulesPath

	rules, err = alert.Store.GetAlertRules()
	if err != nil {
		log.Errorf("Rule File Update Error: %v", err)
		return
	}

	ruleNames = make([]*models.RuleOperation, 0)
	for _, rule := range rules {
		ruleOps := models.NewRuleOperation()
		ruleOps.Rule = rule
		ruleOps.File = fmt.Sprintf("%s-%s.rule", rule.Name, rule.Alert)

		t := template.Must(template.New("ruleTempl").Parse(ruleTempl))
		var buf bytes.Buffer
		err = t.Execute(&buf, rule)
		if err != nil {
			log.Errorf("Rule File Update Error: %v", err)
			return
		}
		h := md5.New()
		io.WriteString(h, buf.String())
		ruleOps.MD5 = h.Sum(nil)
		ruleNames = append(ruleNames, ruleOps)
	}

	files = make(map[string]interface{})
	files, err = getFilelist(path)
	if err != nil {
		log.Errorf("Rule File Update Error: %v", err)
		fmt.Println(err)
		return
	}

	var createRule, deleteRule []*models.RuleOperation
	countRuleNames := len(ruleNames)
	countFiles := len(files)
	if countRuleNames == 0 && countFiles != 0 {
		for k := range files {
			ruleOps := models.NewRuleOperation()
			ruleOps.File = k
			deleteRule = append(deleteRule, ruleOps)
		}
	} else if countFiles == 0 {
		for _, ruleName := range ruleNames {
			createRule = append(createRule, ruleName)
			delete(files, ruleName.File)
		}
	} else {
		for _, ruleName := range ruleNames {
			if files[ruleName.File] == nil {
				createRule = append(createRule, ruleName)
			} else if !bytes.Equal(files[ruleName.File].([]byte), ruleName.MD5) {
				createRule = append(createRule, ruleName)
			}
			delete(files, ruleName.File)
		}
	}

	if len(files) != 0 {
		for k := range files {
			ruleOps := models.NewRuleOperation()
			ruleOps.File = k
			deleteRule = append(deleteRule, ruleOps)
		}
	}

	if len(createRule) != 0 {
		for _, ruleOps := range createRule {
			err := updateFileByAction(path, ruleOps, ruleFileUpdate)
			if err != nil {
				log.Errorf("Rule File Update Error: %v", err)
				return
			}
		}
		reloadReady = true
	}

	if len(deleteRule) != 0 {
		for _, ruleOps := range deleteRule {
			err := updateFileByAction(path, ruleOps, ruleFileDelete)
			if err != nil {
				log.Errorf("Rule File Delete Error: %v", err)
				return
			}
		}
		reloadReady = true
	}

	if reloadReady {
		err := alert.ReloadPrometheusConf()
		if err != nil {
			log.Errorf("Rule File Update Error: %v", err)
			return
		}
		log.Infof("Reload prometheus conf.")
	}

	log.Infof("No rules Changed")
}

func getFilelist(path string) (map[string]interface{}, error) {
	fmt.Println(path)
	files := make(map[string]interface{})
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		md5, _ := getFileMD5value(path)
		files[f.Name()] = md5
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, err
}

func getFileMD5value(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.Write(fileContent)

	h := md5.New()
	io.WriteString(h, buf.String())
	fmt.Println("getFileMD5value", h.Sum(nil))
	return h.Sum(nil), err
}

func updateFileByAction(path string, ruleOps *models.RuleOperation, action string) error {
	file := fmt.Sprintf("%s/%s", path, ruleOps.File)

	if action == ruleFileDelete {
		err := os.Remove(file)
		if err != nil {
			return err
		}
		log.Infof("deleted rule file %s", ruleOps.File)
		return nil
	}

	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return err
	}
	t := template.Must(template.New("ruleTempl").Parse(ruleTempl))
	var buf bytes.Buffer
	err = t.Execute(&buf, ruleOps.Rule)
	if err != nil {
		log.Errorln("executing templta: ", err)
		return err
	}
	f.WriteString(buf.String())
	log.Infof("updated rule file %s", ruleOps.File)
	return nil
}
