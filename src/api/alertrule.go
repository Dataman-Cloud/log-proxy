package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"
	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusrule"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	// ReceiveEventError code
	ReceiveEventError = "503-21000"
	// ReceiveEventError code
	AckEventError = "503-21001"

	// PromtheusReloadPath path string
	PromtheusReloadPath = "/-/reload"

	ruleTempl = `# This rule was update at {{ .UpdatedAt }}
ALERT {{.Alert}}
  IF {{ .Expr }}
  FOR {{ .Pending }}
  LABELS {{ .Labels }}
  ANNOTATIONS {{ .Annotations }}
`
	ruleFileUpdate = "update"
	ruleFileDelete = "delete"
	ruleInterval   = "1m"

	ruleStatusActive   = "Enabled"
	ruleStatusInActive = "Disabled"
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

// GetAlertIndicators return the alert rule indicator list
func (alert *Alert) GetAlertIndicators(ctx *gin.Context) {
	mapper := prometheusrule.NewRuleMapper()
	data := mapper.GetRuleIndicatorsList()
	utils.Ok(ctx, data)
}

// CreateAlertRule create the alert rule in Database
func (alert *Alert) CreateAlertRule(ctx *gin.Context) {
	var (
		data models.Rule
		err  error
	)
	rule := models.NewRule()
	if err = ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	if err = isValidRuleFiled(rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	// Create Alert rule in DB
	err = alert.Store.CreateAlertRule(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	// Get the rule record from DB
	data, err = alert.Store.GetAlertRuleByUniqueIndex(rule.Class, rule.Name, rule.Cluster, rule.App)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	// Write the Rule in file and reload Prometheus conf
	err = alert.WriteAlertFile(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	// Update the rule status as active
	data.Status = ruleStatusActive
	err = alert.Store.UpdateAlertRule(&data)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, data)
}

func isValidRuleFiled(rule *models.Rule) error {
	switch {
	case rule.Class == "":
		return fmt.Errorf("class required")
	case rule.Name == "":
		return fmt.Errorf("name required")
	case rule.Cluster == "":
		return fmt.Errorf("cluster required")
	case rule.App == "":
		return fmt.Errorf("app required")
	case rule.Pending == "":
		return fmt.Errorf("pending required")
	case rule.Indicator == "":
		return fmt.Errorf("indicator required")
	case rule.Severity == "":
		return fmt.Errorf("severity required")
	case rule.Aggregation == "":
		return fmt.Errorf("aggregation required")
	case rule.Comparison == "":
		return fmt.Errorf("comparison required")
	}

	switch {
	case strings.Contains(rule.Class, "_"):
		return fmt.Errorf("Char \"_\" should not exist in class: %s", rule.Class)
	case strings.Contains(rule.Name, "_"):
		return fmt.Errorf("Char \"_\" should not exist in name: %s", rule.Name)
	case strings.Contains(rule.Cluster, "_"):
		return fmt.Errorf("Char \"_\" should not exist in cluster: %s", rule.Cluster)
	case strings.Contains(rule.App, "_"):
		return fmt.Errorf("Char \"_\" should not exist in app: %s", rule.App)
	}

	switch {
	case strings.Contains(rule.Class, "-"):
		return fmt.Errorf("Char \"-\" should not exist in class: %s", rule.Class)
	case strings.Contains(rule.Name, "-"):
		return fmt.Errorf("Char \"-\" should not exist in ame: %s", rule.Name)
	case strings.Contains(rule.Cluster, "-"):
		return fmt.Errorf("Char \"-\" should not exist in cluster: %s", rule.Cluster)
	}

	return nil
}

// ListAlertRules list the rules by name with pages.
func (alert *Alert) ListAlertRules(ctx *gin.Context) {
	class := ctx.Query("class")
	cluster := ctx.Query("cluster")
	app := ctx.Query("app")
	data, err := alert.Store.ListAlertRules(ctx.MustGet("page").(models.Page), class, cluster, app)
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

// DeleteAlertRule delete the rule by id and name from DB and files
func (alert *Alert) DeleteAlertRule(ctx *gin.Context) {
	var (
		rowsAffected int64
		err          error
		result       models.Rule
		id           uint64
		class        string
	)

	id, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, fmt.Errorf("Failed to parse the id: %s", err))
		return
	}
	class = ctx.Query("class")

	// Get alert rule by ID
	result, err = alert.Store.GetAlertRule(id)
	if err != nil {
		log.Errorln("DeleteAlertRule: GetAlertRule() error, ", err)
		utils.ErrorResponse(ctx, err)
		return
	}

	// Delate alert rule
	rowsAffected, err = alert.Store.DeleteAlertRuleByIDClass(id, class)
	if err != nil {
		utils.ErrorResponse(ctx, fmt.Errorf("Failed to delete the rule by %s", err))
		return
	}

	if rowsAffected == 0 {
		utils.Ok(ctx, "no rule was deleted")
		return
	}

	// Update the alert file content
	err = alert.UpdateAlertFile(&result)
	if err != nil {
		log.Errorln("DeleteAlertRule: Delete Alert file error, ", err)
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, "success")
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

	rule.ID, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, fmt.Errorf("Failed to parse the id: %s", err))
		return
	}

	err = alert.Store.UpdateAlertRule(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	data, err = alert.Store.GetAlertRuleByUniqueIndex(rule.Class, rule.Name, rule.Cluster, rule.App)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	if data.Status == ruleStatusInActive {
		err = alert.UpdateAlertFile(&data)
		if err != nil {
			log.Errorln("DeleteAlertRule: Delete Alert file error, ", err)
			utils.ErrorResponse(ctx, err)
			return
		}
	} else {
		err = alert.WriteAlertFile(&data)
		if err != nil {
			utils.ErrorResponse(ctx, err)
			return
		}
	}
	utils.Ok(ctx, data)
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

// WriteAlertFile write the alert rule to file
func (alert *Alert) WriteAlertFile(rule *models.Rule) error {
	var (
		mapper  *prometheusrule.RuleMapper
		rawRule *models.RawRule
		err     error
	)
	// mapping the rule from rule to raw rule
	mapper = prometheusrule.NewRuleMapper()
	rawRule, err = mapper.Map2Raw(rule)
	if err != nil {
		return err
	}
	rawRule.UpdatedAt = time.Now()

	//open the alert rule file
	path := alert.RulesPath
	alertfile := fmt.Sprintf("%s/%s.rule", path, rawRule.Alert)
	f, err := os.Create(alertfile)
	defer f.Close()
	if err != nil {
		return err
	}
	// convert the rawRule with the template
	t := template.Must(template.New("ruleTempl").Parse(ruleTempl))
	var buf bytes.Buffer
	err = t.Execute(&buf, rawRule)
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

// UpdateAlertFile remove the rule from the file.
func (alert *Alert) UpdateAlertFile(rule *models.Rule) error {
	var (
		err     error
		message string
	)

	filename := fmt.Sprintf("%s_%s_%s_%s", rule.Class, rule.Name, rule.Cluster, rule.App)

	path := alert.RulesPath
	alertfile := fmt.Sprintf("%s/%s.rule", path, filename)
	f, err := os.Create(alertfile)
	defer f.Close()
	if err != nil {
		return err
	}

	message = "# inactive this rule"

	f.WriteString(message + "\n")

	err = alert.ReloadPrometheusConf()
	if err != nil {
		return err
	}
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
			AlertName: labels["alertname"].(string),
			Severity:  labels["severity"].(string),
			Cluster:   labels["container_label_VCLUSTER"].(string),
			App:       labels["container_label_APP_ID"].(string),
			Task:      labels["container_env_mesos_task_id"].(string),
			//UserName:    labels["container_label_USER_NAME"].(string),
			//GroupName:   labels["container_label_GROUP_NAME"].(string),
			ContainerID:   labels["id"].(string),
			ContainerName: labels["name"].(string),
			Value:         labels["value"].(string),
			Description:   annotations["description"].(string),
			Summary:       annotations["summary"].(string),
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
		var cluster, app string
		if data["cluster"] != nil {
			cluster = data["cluster"].(string)
		}
		if data["app"] != nil {
			app = data["app"].(string)
		}
		if err = alert.Store.AckEvent(pk, cluster, app); err != nil {
			utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
			return
		}
		utils.Ok(ctx, map[string]string{"status": "success"})
	}
}

// GetAlertEvents list the alert events
func (alert *Alert) GetAlertEvents(ctx *gin.Context) {
	var (
		err error
	)
	options := make(map[string]interface{})
	if ctx.Query("ack") == "" {
		options["ack"] = false
	} else {
		options["ack"], err = strconv.ParseBool(ctx.Query("ack"))
		if err != nil {
			utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
			return
		}
	}

	if ctx.Query("cluster") != "" {
		options["cluster"] = ctx.Query("cluster")
	}
	if ctx.Query("app") != "" {
		options["app"] = ctx.Query("app")
	}
	if ctx.Query("start") != "" {
		options["start"] = ctx.Query("start")
	}
	if ctx.Query("end") != "" {
		options["end"] = ctx.Query("end")
	}

	result, err := alert.Store.ListEvents(ctx.MustGet("page").(models.Page), options)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}
	utils.Ok(ctx, result)
}

/*
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
		fmt.Println(err)
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
*/
