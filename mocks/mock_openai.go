// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jonada182/cover-letter-ai-api/internal/openai (interfaces: OpenAI)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	types "github.com/jonada182/cover-letter-ai-api/types"
	gomock "go.uber.org/mock/gomock"
)

// MockOpenAI is a mock of OpenAI interface.
type MockOpenAI struct {
	ctrl     *gomock.Controller
	recorder *MockOpenAIMockRecorder
}

// MockOpenAIMockRecorder is the mock recorder for MockOpenAI.
type MockOpenAIMockRecorder struct {
	mock *MockOpenAI
}

// NewMockOpenAI creates a new mock instance.
func NewMockOpenAI(ctrl *gomock.Controller) *MockOpenAI {
	mock := &MockOpenAI{ctrl: ctrl}
	mock.recorder = &MockOpenAIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOpenAI) EXPECT() *MockOpenAIMockRecorder {
	return m.recorder
}

// GenerateChatGPTCoverLetter mocks base method.
func (m *MockOpenAI) GenerateChatGPTCoverLetter(arg0 *gin.Context, arg1 string, arg2 *types.JobPosting, arg3 types.StoreClient) (string, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateChatGPTCoverLetter", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateChatGPTCoverLetter indicates an expected call of GenerateChatGPTCoverLetter.
func (mr *MockOpenAIMockRecorder) GenerateChatGPTCoverLetter(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateChatGPTCoverLetter", reflect.TypeOf((*MockOpenAI)(nil).GenerateChatGPTCoverLetter), arg0, arg1, arg2, arg3)
}

// GetCareerProfileInfoPrompt mocks base method.
func (m *MockOpenAI) GetCareerProfileInfoPrompt(arg0 string, arg1 types.StoreClient) (string, *types.CareerProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCareerProfileInfoPrompt", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*types.CareerProfile)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCareerProfileInfoPrompt indicates an expected call of GetCareerProfileInfoPrompt.
func (mr *MockOpenAIMockRecorder) GetCareerProfileInfoPrompt(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCareerProfileInfoPrompt", reflect.TypeOf((*MockOpenAI)(nil).GetCareerProfileInfoPrompt), arg0, arg1)
}

// ParseCoverLetter mocks base method.
func (m *MockOpenAI) ParseCoverLetter(arg0 *string, arg1 *types.CareerProfile, arg2 *types.JobPosting) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseCoverLetter", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseCoverLetter indicates an expected call of ParseCoverLetter.
func (mr *MockOpenAIMockRecorder) ParseCoverLetter(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseCoverLetter", reflect.TypeOf((*MockOpenAI)(nil).ParseCoverLetter), arg0, arg1, arg2)
}
