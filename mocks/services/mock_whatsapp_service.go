// Code generated by MockGen. DO NOT EDIT.
// Source: services/whatsapp_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWhatsappService is a mock of WhatsappService interface.
type MockWhatsappService struct {
	ctrl     *gomock.Controller
	recorder *MockWhatsappServiceMockRecorder
}

// MockWhatsappServiceMockRecorder is the mock recorder for MockWhatsappService.
type MockWhatsappServiceMockRecorder struct {
	mock *MockWhatsappService
}

// NewMockWhatsappService creates a new mock instance.
func NewMockWhatsappService(ctrl *gomock.Controller) *MockWhatsappService {
	mock := &MockWhatsappService{ctrl: ctrl}
	mock.recorder = &MockWhatsappServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWhatsappService) EXPECT() *MockWhatsappServiceMockRecorder {
	return m.recorder
}

// GetMedia mocks base method.
func (m *MockWhatsappService) GetMedia(arg0 string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMedia", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMedia indicates an expected call of GetMedia.
func (mr *MockWhatsappServiceMockRecorder) GetMedia(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMedia", reflect.TypeOf((*MockWhatsappService)(nil).GetMedia), arg0)
}

// Health mocks base method.
func (m *MockWhatsappService) Health() (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health")
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Health indicates an expected call of Health.
func (mr *MockWhatsappServiceMockRecorder) Health() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockWhatsappService)(nil).Health))
}

// Login mocks base method.
func (m *MockWhatsappService) Login() (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login")
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockWhatsappServiceMockRecorder) Login() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockWhatsappService)(nil).Login))
}

// SendMessage mocks base method.
func (m *MockWhatsappService) SendMessage(arg0 []byte) (http.Header, io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0)
	ret0, _ := ret[0].(http.Header)
	ret1, _ := ret[1].(io.ReadCloser)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockWhatsappServiceMockRecorder) SendMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockWhatsappService)(nil).SendMessage), arg0)
}
