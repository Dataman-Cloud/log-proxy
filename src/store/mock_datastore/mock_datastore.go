// Automatically generated by MockGen. DO NOT EDIT!
// Source: ./src/store/store.go

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

func (_m *MockStore) ListAlertRules(page models.Page, user string) (map[string]interface{}, error) {
	ret := _m.ctrl.Call(_m, "ListAlertRules", page, user)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) ListAlertRules(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListAlertRules", arg0, arg1)
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

func (_m *MockStore) GetAlertRuleByName(name string, alert string) (models.Rule, error) {
	ret := _m.ctrl.Call(_m, "GetAlertRuleByName", name, alert)
	ret0, _ := ret[0].(models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) GetAlertRuleByName(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAlertRuleByName", arg0, arg1)
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

func (_m *MockStore) DeleteAlertRuleByIDName(id uint64, name string) (int64, error) {
	ret := _m.ctrl.Call(_m, "DeleteAlertRuleByIDName", id, name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) DeleteAlertRuleByIDName(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteAlertRuleByIDName", arg0, arg1)
}

func (_m *MockStore) CreateOrIncreaseEvent(event *models.Event) error {
	ret := _m.ctrl.Call(_m, "CreateOrIncreaseEvent", event)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) CreateOrIncreaseEvent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateOrIncreaseEvent", arg0)
}

func (_m *MockStore) AckEvent(pk int, username string, groupname string) error {
	ret := _m.ctrl.Call(_m, "AckEvent", pk, username, groupname)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) AckEvent(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AckEvent", arg0, arg1, arg2)
}

func (_m *MockStore) ListAckedEvent(page models.Page, username string, groupname string) map[string]interface{} {
	ret := _m.ctrl.Call(_m, "ListAckedEvent", page, username, groupname)
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

func (_mr *_MockStoreRecorder) ListAckedEvent(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListAckedEvent", arg0, arg1, arg2)
}

func (_m *MockStore) ListUnackedEvent(page models.Page, username string, groupname string) map[string]interface{} {
	ret := _m.ctrl.Call(_m, "ListUnackedEvent", page, username, groupname)
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

func (_mr *_MockStoreRecorder) ListUnackedEvent(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListUnackedEvent", arg0, arg1, arg2)
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
