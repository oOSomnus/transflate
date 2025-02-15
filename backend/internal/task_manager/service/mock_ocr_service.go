// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/task_manager/service/ocr_service.go

// Package mocks is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ocr "github.com/oOSomnus/transflate/api/generated/ocr"
)

// MockOCRClient is a mock of OCRClient interface.
type MockOCRClient struct {
	ctrl     *gomock.Controller
	recorder *MockOCRClientMockRecorder
}

// MockOCRClientMockRecorder is the mock recorder for MockOCRClient.
type MockOCRClientMockRecorder struct {
	mock *MockOCRClient
}

// NewMockOCRClient creates a new mock instance.
func NewMockOCRClient(ctrl *gomock.Controller) *MockOCRClient {
	mock := &MockOCRClient{ctrl: ctrl}
	mock.recorder = &MockOCRClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOCRClient) EXPECT() *MockOCRClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockOCRClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockOCRClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockOCRClient)(nil).Close))
}

// ProcessOCR mocks base method.
func (m *MockOCRClient) ProcessOCR(fileContent []byte, lang string) (*ocr.StringListResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessOCR", fileContent, lang)
	ret0, _ := ret[0].(*ocr.StringListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessOCR indicates an expected call of ProcessOCR.
func (mr *MockOCRClientMockRecorder) ProcessOCR(fileContent, lang interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessOCR", reflect.TypeOf((*MockOCRClient)(nil).ProcessOCR), fileContent, lang)
}
