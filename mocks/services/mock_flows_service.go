// Code generated by MockGen. DO NOT EDIT.
// Source: services/flows_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
        reflect "reflect"

        gomock "github.com/golang/mock/gomock"
        models "github.com/weni/whatsapp-router/models"
)

// MockFlowsService is a mock of FlowsService interface.
type MockFlowsService struct {
        ctrl     *gomock.Controller
        recorder *MockFlowsServiceMockRecorder
}

// MockFlowsServiceMockRecorder is the mock recorder for MockFlowsService.
type MockFlowsServiceMockRecorder struct {
        mock *MockFlowsService
}

// NewMockFlowsService creates a new mock instance.
func NewMockFlowsService(ctrl *gomock.Controller) *MockFlowsService {
        mock := &MockFlowsService{ctrl: ctrl}
        mock.recorder = &MockFlowsServiceMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlowsService) EXPECT() *MockFlowsServiceMockRecorder {
        return m.recorder
}

// CreateFlows mocks base method.
func (m *MockFlowsService) CreateFlows(arg0 *models.Flows) (*models.Flows, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "CreateFlows", arg0)
        ret0, _ := ret[0].(*models.Flows)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// CreateFlows indicates an expected call of CreateFlows.
func (mr *MockFlowsServiceMockRecorder) CreateFlows(arg0 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFlows", reflect.TypeOf((*MockFlowsService)(nil).CreateFlows), arg0)
}

// FindFlows mocks base method.
func (m *MockFlowsService) FindFlows(arg0 *models.Flows) (*models.Flows, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindFlows", arg0)
        ret0, _ := ret[0].(*models.Flows)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindFlows indicates an expected call of FindFlows.
func (mr *MockFlowsServiceMockRecorder) FindFlows(arg0 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFlows", reflect.TypeOf((*MockFlowsService)(nil).FindFlows), arg0)
}

// UpdateFlows mocks base method.
func (m *MockFlowsService) UpdateFlows(arg0 *models.Flows) (*models.Flows, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "UpdateFlows", arg0)
        ret0, _ := ret[0].(*models.Flows)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// UpdateFlows indicates an expected call of UpdateFlows.
func (mr *MockFlowsServiceMockRecorder) UpdateFlows(arg0 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlows", reflect.TypeOf((*MockFlowsService)(nil).UpdateFlows), arg0)
}