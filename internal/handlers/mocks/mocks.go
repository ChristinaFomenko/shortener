// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Mockservice is a mock of service interface.
type Mockservice struct {
	ctrl     *gomock.Controller
	recorder *MockserviceMockRecorder
}

// MockserviceMockRecorder is the mock recorder for Mockservice.
type MockserviceMockRecorder struct {
	mock *Mockservice
}

// NewMockservice creates a new mock instance.
func NewMockservice(ctrl *gomock.Controller) *Mockservice {
	mock := &Mockservice{ctrl: ctrl}
	mock.recorder = &MockserviceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockservice) EXPECT() *MockserviceMockRecorder {
	return m.recorder
}

// Expand mocks base method.
func (m *Mockservice) Expand(id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expand", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Expand indicates an expected call of Expand.
func (mr *MockserviceMockRecorder) Expand(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expand", reflect.TypeOf((*Mockservice)(nil).Expand), id)
}

// GetByUsers mocks base method.
func (m *Mockservice) GetByUsers(UserID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsers", UserID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsers indicates an expected call of GetByUsers.
func (mr *MockserviceMockRecorder) GetByUsers(UserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsers", reflect.TypeOf((*Mockservice)(nil).GetByUsers), UserID)
}

// Shorten mocks base method.
func (m *Mockservice) Shorten(url string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shorten", url)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Shorten indicates an expected call of Shorten.
func (mr *MockserviceMockRecorder) Shorten(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shorten", reflect.TypeOf((*Mockservice)(nil).Shorten), url)
}
