package alertservice

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/router/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_store "github.com/Dataman-Cloud/log-proxy/src/store/mock_datastore"
)

func TestNewAlert(t *testing.T) {
	alert := NewAlert()
	if alert == nil {
		t.Errorf("expect not nil, got %v", alert)
	}
}

func TestGetAlertIndicators(t *testing.T) {
	alert := NewAlert()
	result := alert.GetAlertIndicators()
	if result == nil {
		t.Errorf("expect not nil, got %v", result)
	}
}

func TestCreateAlertRule(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	var result models.Rule
	result = *rule
	result.ID = 1

	mockStore.EXPECT().CreateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRuleByName(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(nil).Times(1)

	data, err := alert.CreateAlertRule(rule)
	if data.Status != "Enabled" {
		t.Errorf("Expect data.Status is Enabled, but got %s", data.Status)
	}
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}

func TestCreateAlertRuleError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.App = "work-nginx"

	var result models.Rule
	result = *rule
	result.ID = 1

	data, err := alert.CreateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is nil, but got %s", data.Status)
	}
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}
	// test wrong indicator
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "error"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	data, err = alert.CreateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is nil, but got %s", data.Status)
	}
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}
	// test create error
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	result = *rule
	result.ID = 1

	mockStore.EXPECT().CreateAlertRule(gomock.Any()).Return(errors.New("error")).Times(1)
	data, err = alert.CreateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is nil, but got %s", data.Status)
	}
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	// test Get error
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	result = *rule
	result.ID = 1

	mockStore.EXPECT().CreateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRuleByName(gomock.Any()).Return(result, errors.New("error")).Times(1)
	data, err = alert.CreateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is nil, but got %s", data.Status)
	}
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	// test Update error
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	result = *rule
	result.ID = 1

	mockStore.EXPECT().CreateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRuleByName(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(errors.New("error")).Times(1)
	data, err = alert.CreateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is nil, but got %s", data.Status)
	}
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}
}

func TestIsValidRuleFile(t *testing.T) {
	var err error

	rule := models.NewRule()
	err = isValidRuleFile(rule)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	rule.App = "app"
	err = isValidRuleFile(rule)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	rule.App = "app"
	rule.Pending = "5s"
	rule.Severity = ""
	err = isValidRuleFile(rule)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	rule.App = "app"
	rule.Pending = "5s"
	rule.Severity = "warning"
	err = isValidRuleFile(rule)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	rule.App = "app"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	err = isValidRuleFile(rule)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	rule.App = "app"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	err = isValidRuleFile(rule)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	rule.App = "app"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	rule.Comparison = ""
	err = isValidRuleFile(rule)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}
}

func TestDeleteAlertRule(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	var result models.Rule
	result = *rule
	result.ID = 1

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().DeleteAlertRuleByID(gomock.Any()).Return(int64(1), nil).Times(1)

	err := alert.DeleteAlertRule(result.ID, "Undefine")
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}

func TestDeleteAlertRuleError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "cpu_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	var result models.Rule
	result = *rule
	result.ID = 1

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, errors.New("error")).Times(1)
	err := alert.DeleteAlertRule(result.ID, "Undefine")
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	err = alert.DeleteAlertRule(result.ID, "error_group")
	if err == nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().DeleteAlertRuleByID(gomock.Any()).Return(int64(1), errors.New("error")).Times(1)
	err = alert.DeleteAlertRule(result.ID, "Undefine")
	if err == nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().DeleteAlertRuleByID(gomock.Any()).Return(int64(0), nil).Times(1)
	err = alert.DeleteAlertRule(result.ID, "Undefine")
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}

func TestListAlertRules(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.Group = "work"
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "mem_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	var result *models.RulesList
	result = models.NewRulesList()
	result.Count = 1
	result.Rules = append(result.Rules, rule)

	mockStore.EXPECT().ListAlertRules(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil).Times(1)

	page := &models.Page{}
	groups := []string{}
	data, err := alert.ListAlertRules(*page, groups, "")
	if data == nil {
		t.Errorf("Expect data is not nil, bug got %v", data)
	}
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}

func TestGetAlertRule(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "mem_usage"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	var result models.Rule
	result = *rule

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)

	data, err := alert.GetAlertRule(uint64(1))
	if data == nil {
		t.Errorf("Expect data is not nil, bug got %v", data)
	}
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}

func TestUpdateAlertRule(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.ID = 1
	rule.Pending = "5s"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60
	rule.Group = "Undefine"
	rule.Status = "Enabled"

	var result models.Rule
	result = *rule
	result.Indicator = "mem_usage"

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)

	data, err := alert.UpdateAlertRule(rule)
	if data == nil {
		t.Errorf("Expect data is not nil, bug got %v", data)
	}
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}

func TestUpdateAlertRuleDisabled(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.ID = 1
	rule.Pending = "5s"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60
	rule.Group = "Undefine"
	rule.Status = "Disabled"

	var result models.Rule
	result = *rule
	result.Indicator = "mem_usage"

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)

	data, err := alert.UpdateAlertRule(rule)
	if data == nil {
		t.Errorf("Expect data is not nil, bug got %v", data)
	}
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}

func TestUpdateAlertRuleError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	rule := models.NewRule()
	rule.ID = 1
	rule.Pending = "5s"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60
	rule.Group = "Undefine"
	rule.Status = "Enabled"
	rule.App = "app"

	var result models.Rule
	result = *rule
	result.Indicator = "mem_usage"

	data, err := alert.UpdateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is not nil, bug got %v", data)
	}
	if err == nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}

	rule.App = ""
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, errors.New("err")).Times(1)
	data, err = alert.UpdateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is not nil, bug got %v", data)
	}
	if err == nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}

	rule.Group = "errgroup"
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	data, err = alert.UpdateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is not nil, bug got %v", data)
	}
	if err == nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}

	rule.Group = "Undefine"
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(errors.New("err")).Times(1)
	data, err = alert.UpdateAlertRule(rule)
	if data != nil {
		t.Errorf("Expect data is not nil, bug got %v", data)
	}
	if err == nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}
