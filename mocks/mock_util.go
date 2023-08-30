// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jonada182/cover-letter-ai-api/util (interfaces: Util)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUtil is a mock of Util interface.
type MockUtil struct {
	ctrl     *gomock.Controller
	recorder *MockUtilMockRecorder
}

// MockUtilMockRecorder is the mock recorder for MockUtil.
type MockUtilMockRecorder struct {
	mock *MockUtil
}

// NewMockUtil creates a new mock instance.
func NewMockUtil(ctrl *gomock.Controller) *MockUtil {
	mock := &MockUtil{ctrl: ctrl}
	mock.recorder = &MockUtilMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUtil) EXPECT() *MockUtilMockRecorder {
	return m.recorder
}

// LoadEnvFile mocks base method.
func (m *MockUtil) LoadEnvFile(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadEnvFile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadEnvFile indicates an expected call of LoadEnvFile.
func (mr *MockUtilMockRecorder) LoadEnvFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadEnvFile", reflect.TypeOf((*MockUtil)(nil).LoadEnvFile), arg0)
}
