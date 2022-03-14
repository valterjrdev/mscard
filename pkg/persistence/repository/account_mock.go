// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/persistence/repository/account.go

// Package repository is a generated GoMock package.
package repository

import (
	entity "ms/card/pkg/persistence/entity"
	filter "ms/card/pkg/persistence/filter"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	context "golang.org/x/net/context"
)

// MockAccounts is a mock of Accounts interface.
type MockAccounts struct {
	ctrl     *gomock.Controller
	recorder *MockAccountsMockRecorder
}

// MockAccountsMockRecorder is the mock recorder for MockAccounts.
type MockAccountsMockRecorder struct {
	mock *MockAccounts
}

// NewMockAccounts creates a new mock instance.
func NewMockAccounts(ctrl *gomock.Controller) *MockAccounts {
	mock := &MockAccounts{ctrl: ctrl}
	mock.recorder = &MockAccountsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccounts) EXPECT() *MockAccountsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAccounts) Create(ctx context.Context, structure entity.Account) (*entity.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, structure)
	ret0, _ := ret[0].(*entity.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAccountsMockRecorder) Create(ctx, structure interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAccounts)(nil).Create), ctx, structure)
}

// FindAll mocks base method.
func (m *MockAccounts) FindAll(ctx context.Context, filters filter.AccountCollection) ([]*entity.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, filters)
	ret0, _ := ret[0].([]*entity.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockAccountsMockRecorder) FindAll(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockAccounts)(nil).FindAll), ctx, filters)
}

// FindByID mocks base method.
func (m *MockAccounts) FindByID(ctx context.Context, id uint) (*entity.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*entity.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockAccountsMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockAccounts)(nil).FindByID), ctx, id)
}