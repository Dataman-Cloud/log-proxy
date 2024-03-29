// Automatically generated by MockGen. DO NOT EDIT!
// Source: store.go

package mock_store

import (
	models "github.com/Dataman-Cloud/log-proxy/src/models"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *_MockStoreRecorder
}

// Recorder for MockStore (not exported)
type _MockStoreRecorder struct {
	mock *MockStore
}

func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &_MockStoreRecorder{mock}
	return mock
}

func (_m *MockStore) EXPECT() *_MockStoreRecorder {
	return _m.recorder
}

func (_m *MockStore) ListAlertRules(page models.Page, groups []string, app string) (*models.RulesList, error) {
	ret := _m.ctrl.Call(_m, "ListAlertRules", page, groups, app)
	ret0, _ := ret[0].(*models.RulesList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) ListAlertRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListAlertRules", arg0, arg1, arg2)
}

func (_m *MockStore) GetAlertRule(id uint64) (models.Rule, error) {
	ret := _m.ctrl.Call(_m, "GetAlertRule", id)
	ret0, _ := ret[0].(models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAlertRule", arg0)
}

func (_m *MockStore) GetAlertRules() ([]*models.Rule, error) {
	ret := _m.ctrl.Call(_m, "GetAlertRules")
	ret0, _ := ret[0].([]*models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetAlertRules() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAlertRules")
}

func (_m *MockStore) GetAlertRuleByName(name string) (models.Rule, error) {
	ret := _m.ctrl.Call(_m, "GetAlertRuleByName", name)
	ret0, _ := ret[0].(models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetAlertRuleByName(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAlertRuleByName", arg0)
}

func (_m *MockStore) CreateAlertRule(rule *models.Rule) error {
	ret := _m.ctrl.Call(_m, "CreateAlertRule", rule)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) CreateAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateAlertRule", arg0)
}

func (_m *MockStore) UpdateAlertRule(rule *models.Rule) error {
	ret := _m.ctrl.Call(_m, "UpdateAlertRule", rule)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) UpdateAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateAlertRule", arg0)
}

func (_m *MockStore) DeleteAlertRuleByID(id uint64) (int64, error) {
	ret := _m.ctrl.Call(_m, "DeleteAlertRuleByID", id)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) DeleteAlertRuleByID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteAlertRuleByID", arg0)
}

func (_m *MockStore) CreateOrIncreaseEvent(event *models.Event) error {
	ret := _m.ctrl.Call(_m, "CreateOrIncreaseEvent", event)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) CreateOrIncreaseEvent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateOrIncreaseEvent", arg0)
}

func (_m *MockStore) AckEvent(ID int, group string, app string) error {
	ret := _m.ctrl.Call(_m, "AckEvent", ID, group, app)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) AckEvent(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AckEvent", arg0, arg1, arg2)
}

func (_m *MockStore) ListEvents(page models.Page, options map[string]interface{}, groups []string) (map[string]interface{}, error) {
	ret := _m.ctrl.Call(_m, "ListEvents", page, options, groups)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) ListEvents(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListEvents", arg0, arg1, arg2)
}

func (_m *MockStore) CreateLogAlertRule(rule *models.LogAlertRule) error {
	ret := _m.ctrl.Call(_m, "CreateLogAlertRule", rule)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) CreateLogAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLogAlertRule", arg0)
}

func (_m *MockStore) UpdateLogAlertRule(rule *models.LogAlertRule) error {
	ret := _m.ctrl.Call(_m, "UpdateLogAlertRule", rule)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) UpdateLogAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateLogAlertRule", arg0)
}

func (_m *MockStore) DeleteLogAlertRule(ID string) error {
	ret := _m.ctrl.Call(_m, "DeleteLogAlertRule", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) DeleteLogAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteLogAlertRule", arg0)
}

func (_m *MockStore) GetLogAlertRule(ID string) (models.LogAlertRule, error) {
	ret := _m.ctrl.Call(_m, "GetLogAlertRule", ID)
	ret0, _ := ret[0].(models.LogAlertRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetLogAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLogAlertRule", arg0)
}

func (_m *MockStore) GetLogAlertRules(opts map[string]interface{}, page models.Page) (map[string]interface{}, error) {
	ret := _m.ctrl.Call(_m, "GetLogAlertRules", opts, page)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetLogAlertRules(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLogAlertRules", arg0, arg1)
}

func (_m *MockStore) GetLogAlertApps(opts map[string]interface{}, page models.Page) ([]*models.LogAlertApps, error) {
	ret := _m.ctrl.Call(_m, "GetLogAlertApps", opts, page)
	ret0, _ := ret[0].([]*models.LogAlertApps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetLogAlertApps(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLogAlertApps", arg0, arg1)
}

func (_m *MockStore) AckLogAlertEvent(ID string) error {
	ret := _m.ctrl.Call(_m, "AckLogAlertEvent", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) AckLogAlertEvent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AckLogAlertEvent", arg0)
}

func (_m *MockStore) GetLogAlertEvent(ID string) (*models.LogAlertEvent, error) {
	ret := _m.ctrl.Call(_m, "GetLogAlertEvent", ID)
	ret0, _ := ret[0].(*models.LogAlertEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetLogAlertEvent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLogAlertEvent", arg0)
}

func (_m *MockStore) GetLogAlertEvents(opts map[string]interface{}, page models.Page) (map[string]interface{}, error) {
	ret := _m.ctrl.Call(_m, "GetLogAlertEvents", opts, page)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetLogAlertEvents(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLogAlertEvents", arg0, arg1)
}

func (_m *MockStore) CreateLogAlertEvent(event *models.LogAlertEvent) error {
	ret := _m.ctrl.Call(_m, "CreateLogAlertEvent", event)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) CreateLogAlertEvent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLogAlertEvent", arg0)
}
