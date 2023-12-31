// Code generated by MockGen. DO NOT EDIT.
// Source: poktroll/x/servicer/keeper (interfaces: ServicerKeeper)

// Package mocks is a generated GoMock package.
package mocks

import (
	types0 "poktroll/x/servicer/types"
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

// MockServicerKeeper is a mock of ServicerKeeper interface.
type MockServicerKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockServicerKeeperMockRecorder
}

// MockServicerKeeperMockRecorder is the mock recorder for MockServicerKeeper.
type MockServicerKeeperMockRecorder struct {
	mock *MockServicerKeeper
}

// NewMockServicerKeeper creates a new mock instance.
func NewMockServicerKeeper(ctrl *gomock.Controller) *MockServicerKeeper {
	mock := &MockServicerKeeper{ctrl: ctrl}
	mock.recorder = &MockServicerKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServicerKeeper) EXPECT() *MockServicerKeeperMockRecorder {
	return m.recorder
}

// GetAllServicers mocks base method.
func (m *MockServicerKeeper) GetAllServicers(arg0 types.Context) []types0.Servicers {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllServicers", arg0)
	ret0, _ := ret[0].([]types0.Servicers)
	return ret0
}

// GetAllServicers indicates an expected call of GetAllServicers.
func (mr *MockServicerKeeperMockRecorder) GetAllServicers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllServicers", reflect.TypeOf((*MockServicerKeeper)(nil).GetAllServicers), arg0)
}

// GetServicers mocks base method.
func (m *MockServicerKeeper) GetServicers(arg0 types.Context, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetServicers", arg0, arg1)
}

// GetServicers indicates an expected call of GetServicers.
func (mr *MockServicerKeeperMockRecorder) GetServicers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServicers", reflect.TypeOf((*MockServicerKeeper)(nil).GetServicers), arg0, arg1)
}

// RemoveServicers mocks base method.
func (m *MockServicerKeeper) RemoveServicers(arg0 types.Context, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveServicers", arg0, arg1)
}

// RemoveServicers indicates an expected call of RemoveServicers.
func (mr *MockServicerKeeperMockRecorder) RemoveServicers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServicers", reflect.TypeOf((*MockServicerKeeper)(nil).RemoveServicers), arg0, arg1)
}

// SetServicers mocks base method.
func (m *MockServicerKeeper) SetServicers(arg0 types.Context, arg1 types0.Servicers) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetServicers", arg0, arg1)
}

// SetServicers indicates an expected call of SetServicers.
func (mr *MockServicerKeeperMockRecorder) SetServicers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetServicers", reflect.TypeOf((*MockServicerKeeper)(nil).SetServicers), arg0, arg1)
}
