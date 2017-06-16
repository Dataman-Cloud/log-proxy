// Automatically generated by MockGen. DO NOT EDIT!
// Source: interface_alert.go

package mock_service

import (
	models "github.com/Dataman-Cloud/log-proxy/src/models"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Alerter interface
type MockAlerter struct {
	ctrl     *gomock.Controller
	recorder *_MockAlerterRecorder
}

// Recorder for MockAlerter (not exported)
type _MockAlerterRecorder struct {
	mock *MockAlerter
}

func NewMockAlerter(ctrl *gomock.Controller) *MockAlerter {
	mock := &MockAlerter{ctrl: ctrl}
	mock.recorder = &_MockAlerterRecorder{mock}
	return mock
}

func (_m *MockAlerter) EXPECT() *_MockAlerterRecorder {
	return _m.recorder
}

func (_m *MockAlerter) GetAlertIndicators() map[string]string {
	ret := _m.ctrl.Call(_m, "GetAlertIndicators")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

func (_mr *_MockAlerterRecorder) GetAlertIndicators() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAlertIndicators")
}

func (_m *MockAlerter) CreateAlertRule(rule *models.Rule) (*models.Rule, error) {
	ret := _m.ctrl.Call(_m, "CreateAlertRule", rule)
	ret0, _ := ret[0].(*models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAlerterRecorder) CreateAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateAlertRule", arg0)
}

func (_m *MockAlerter) GetIndicatorName(alias string) (string, string, error) {
	ret := _m.ctrl.Call(_m, "GetIndicatorName", alias)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockAlerterRecorder) GetIndicatorName(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetIndicatorName", arg0)
}

func (_m *MockAlerter) GetIndicatorAlias(name string) (string, string, error) {
	ret := _m.ctrl.Call(_m, "GetIndicatorAlias", name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockAlerterRecorder) GetIndicatorAlias(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetIndicatorAlias", arg0)
}

func (_m *MockAlerter) ReloadPrometheusConf() error {
	ret := _m.ctrl.Call(_m, "ReloadPrometheusConf")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAlerterRecorder) ReloadPrometheusConf() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReloadPrometheusConf")
}

func (_m *MockAlerter) WriteAlertFile(rule *models.Rule) error {
	ret := _m.ctrl.Call(_m, "WriteAlertFile", rule)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAlerterRecorder) WriteAlertFile(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WriteAlertFile", arg0)
}

func (_m *MockAlerter) DeleteAlertRule(id uint64, groups []string) error {
	ret := _m.ctrl.Call(_m, "DeleteAlertRule", id, groups)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAlerterRecorder) DeleteAlertRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteAlertRule", arg0, arg1)
}

func (_m *MockAlerter) UpdateAlertFile(rule *models.Rule) error {
	ret := _m.ctrl.Call(_m, "UpdateAlertFile", rule)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAlerterRecorder) UpdateAlertFile(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateAlertFile", arg0)
}

func (_m *MockAlerter) ListAlertRules(page models.Page, groups []string, app string) (*models.RulesList, error) {
	ret := _m.ctrl.Call(_m, "ListAlertRules", page, groups, app)
	ret0, _ := ret[0].(*models.RulesList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAlerterRecorder) ListAlertRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListAlertRules", arg0, arg1, arg2)
}

func (_m *MockAlerter) GetAlertRule(id uint64) (*models.Rule, error) {
	ret := _m.ctrl.Call(_m, "GetAlertRule", id)
	ret0, _ := ret[0].(*models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAlerterRecorder) GetAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAlertRule", arg0)
}

func (_m *MockAlerter) UpdateAlertRule(rule *models.Rule) (*models.Rule, error) {
	ret := _m.ctrl.Call(_m, "UpdateAlertRule", rule)
	ret0, _ := ret[0].(*models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAlerterRecorder) UpdateAlertRule(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateAlertRule", arg0)
}

func (_m *MockAlerter) ReceiveAlertEvent(message map[string]interface{}) error {
	ret := _m.ctrl.Call(_m, "ReceiveAlertEvent", message)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAlerterRecorder) ReceiveAlertEvent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReceiveAlertEvent", arg0)
}

func (_m *MockAlerter) GetAlertEvents(page models.Page, options map[string]interface{}, groups []string) (map[string]interface{}, error) {
	ret := _m.ctrl.Call(_m, "GetAlertEvents", page, options, groups)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAlerterRecorder) GetAlertEvents(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAlertEvents", arg0, arg1, arg2)
}

func (_m *MockAlerter) AckAlertEvent(id int, options map[string]interface{}) error {
	ret := _m.ctrl.Call(_m, "AckAlertEvent", id, options)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAlerterRecorder) AckAlertEvent(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AckAlertEvent", arg0, arg1)
}
