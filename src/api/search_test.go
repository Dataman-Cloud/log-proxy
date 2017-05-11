package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	mock_service "github.com/Dataman-Cloud/log-proxy/src/service/mock_log_search"
	mock_store "github.com/Dataman-Cloud/log-proxy/src/store/mock_datastore"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetLogAlertRuleIndex(t *testing.T) {
	rule := models.LogAlertRule{
		Group:  "g",
		User:   "u",
		App:    "a",
		Source: "s",
	}

	ruleIndex := getLogAlertRuleIndex(rule)
	assert.Equal(t, ruleIndex, "g-a-s")
}

func TestInitLogAlertFilter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := GetSearch()
	s.Store = mockStore

	rule := models.LogAlertRule{
		App:     "app",
		Cluster: "cluster",
		Keyword: "key",
		Source:  "stdout",
		User:    "user",
		Group:   "group",
	}

	rules := []*models.LogAlertRule{&rule}
	result := map[string]interface{}{"rules": rules}
	mockStore.EXPECT().GetLogAlertRules(gomock.Any(), gomock.Any()).Return(result, nil).Times(1)
	s.InitLogKeywordFilter()

	mockStore.EXPECT().GetLogAlertRules(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	s.InitLogKeywordFilter()

	mockStore.EXPECT().GetLogAlertRules(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	s.InitLogKeywordFilter()
}

func TestMain(m *testing.M) {
	config.InitConfig("../../env_file.template")
	config.LoadLogOptionalLabels()
	ret := m.Run()
	os.Exit(ret)
}

func TestPing(t *testing.T) {
	router := gin.New()
	s := Search{}
	router.GET("/v1/ping", s.Ping)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/v1/ping")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestApplications(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockService := mock_service.NewMockLogSearchService(mockCtl)
	s := Search{Service: mockService}
	router := gin.New()
	router.GET("/v1/apps", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.Applications)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockService.EXPECT().Applications(gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err := http.Get(testServer.URL + "/v1/apps")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	mockService.EXPECT().Applications(gomock.Any()).Return(nil, nil).Times(1)
	resp, err = http.Get(testServer.URL + "/v1/apps")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestSlots(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockService := mock_service.NewMockLogSearchService(mockCtl)
	s := Search{Service: mockService}
	router := gin.New()
	router.GET("/v1/apps/test/slots", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.Slots)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockService.EXPECT().Slots(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err := http.Get(testServer.URL + "/v1/apps/test/slots")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	mockService.EXPECT().Slots(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	resp, err = http.Get(testServer.URL + "/v1/apps/test/slots")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestTasks(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockService := mock_service.NewMockLogSearchService(mockCtl)
	s := Search{Service: mockService}
	router := gin.New()
	router.GET("/v1/apps/test/tasks", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.Tasks)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockService.EXPECT().Tasks(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err := http.Get(testServer.URL + "/v1/apps/test/tasks?slot=0")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	mockService.EXPECT().Tasks(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	resp, err = http.Get(testServer.URL + "/v1/apps/test/tasks?slot=0")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestSources(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockService := mock_service.NewMockLogSearchService(mockCtl)
	s := Search{Service: mockService}
	router := gin.New()
	router.GET("/v1/apps/test/sources", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.Sources)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockService.EXPECT().Sources(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err := http.Get(testServer.URL + "/v1/apps/test/sources?slot=0&task=test")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	mockService.EXPECT().Sources(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	resp, err = http.Get(testServer.URL + "/v1/apps/test/sources?slot=0&task=test")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestSearch(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockService := mock_service.NewMockLogSearchService(mockCtl)
	s := Search{Service: mockService}
	router := gin.New()
	router.GET("/v1/apps/test/search", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.Search)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockService.EXPECT().Search(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err := http.Get(testServer.URL + "/v1/apps/test/search?slot=0&task=test")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	mockService.EXPECT().Search(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	resp, err = http.Get(testServer.URL + "/v1/apps/test/search?slot=0&task=test")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestContext(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockService := mock_service.NewMockLogSearchService(mockCtl)
	s := Search{Service: mockService}
	router := gin.New()
	router.GET("/v1/apps/test/context", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.Context)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockService.EXPECT().Context(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err := http.Get(testServer.URL + "/v1/apps/test/context?slot=0&task=test")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	mockService.EXPECT().Context(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	resp, err = http.Get(testServer.URL + "/v1/apps/test/context?slot=0&task=test")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestEverything(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockService := mock_service.NewMockLogSearchService(mockCtl)
	s := Search{Service: mockService}
	router := gin.New()
	router.GET("/v1/everything/app", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.Everything)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockService.EXPECT().Everything(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err := http.Get(testServer.URL + "/v1/everything/app?slot=0&task=test")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	mockService.EXPECT().Everything(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	resp, err = http.Get(testServer.URL + "/v1/everything/app?slot=0&task=test")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
