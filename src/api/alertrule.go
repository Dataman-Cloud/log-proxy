package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	// ReceiveEventError code
	ReceiveEventError = "503-21000"
	// AckEventError code
	AckEventError = "503-21001"
)

// GetAlertIndicators return the alert rule indicator list
func (m *Monitor) GetAlertIndicators(ctx *gin.Context) {
	utils.Ok(ctx, m.Alert.GetAlertIndicators())
}

// CreateAlertRule create the alert rule in Database
func (m *Monitor) CreateAlertRule(ctx *gin.Context) {
	var (
		data *models.Rule
		err  error
	)
	rule := models.NewRule()
	if err = ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	data, err = m.Alert.CreateAlertRule(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// DeleteAlertRule delete the rule by id and name from DB and files
func (m *Monitor) DeleteAlertRule(ctx *gin.Context) {
	var (
		err   error
		id    uint64
		group string
	)

	id, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, fmt.Errorf("Failed to parse the id: %s", err))
		return
	}
	group = ctx.Query("group")

	err = m.Alert.DeleteAlertRule(id, group)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, "success")
}

// ListAlertRules list the rules by name with pages.
func (m *Monitor) ListAlertRules(ctx *gin.Context) {
	page := ctx.MustGet("page").(models.Page)
	groups := ctx.QueryArray("group")
	app := ctx.Query("app")
	data, err := m.Alert.ListAlertRules(page, groups, app)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, data)
}

// GetAlertRule return the info of alert rule by id
func (m *Monitor) GetAlertRule(ctx *gin.Context) {
	var (
		data *models.Rule
		err  error
		id   uint64
	)

	id, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	data, err = m.Alert.GetAlertRule(id)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, data)
}

// UpdateAlertRule update the alert rule in Database
func (m *Monitor) UpdateAlertRule(ctx *gin.Context) {
	var (
		rule, data *models.Rule
		err        error
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

	data, err = m.Alert.UpdateAlertRule(rule)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// ReceiveAlertEvent recive the alerts from Alertmanager
func (m *Monitor) ReceiveAlertEvent(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ReceiveEventError, err))
		return
	}
	log.Infof("Got event from alertmanager: %s", string(data))
	var message map[string]interface{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ReceiveEventError, err))
		return
	}

	err = m.Alert.ReceiveAlertEvent(message)
	if err != nil {
		utils.ErrorResponse(ctx, fmt.Errorf("Receive the alert message with err %v", err))
		return
	}
	utils.Ok(ctx, map[string]string{"status": "success"})
}

// GetAlertEvents list the alert events
func (m *Monitor) GetAlertEvents(ctx *gin.Context) {
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

	if ctx.Query("group") != "" {
		options["groupname"] = ctx.Query("group")
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

	result, err := m.Alert.GetAlertEvents(ctx.MustGet("page").(models.Page), options)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}
	utils.Ok(ctx, result)
}

// AckAlertEvent mark the alert evnet ACK
func (m *Monitor) AckAlertEvent(ctx *gin.Context) {
	var (
		err error
	)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}
	var options map[string]interface{}

	if err = ctx.BindJSON(&options); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}

	err = m.Alert.AckAlertEvent(id, options)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}

	utils.Ok(ctx, map[string]string{"status": "success"})
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
