// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store/template.go
//
// Generated by this command:
//
//	mockgen -source=internal/store/template.go -destination internal/store/mock/template_mock.go
//
// Package mock_store is a generated GoMock package.
package mock_store

import (
	context "context"
	reflect "reflect"

	models "github.com/emergency-messages/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockTemplater is a mock of Templater interface.
type MockTemplater struct {
	ctrl     *gomock.Controller
	recorder *MockTemplaterMockRecorder
}

// MockTemplaterMockRecorder is the mock recorder for MockTemplater.
type MockTemplaterMockRecorder struct {
	mock *MockTemplater
}

// NewMockTemplater creates a new mock instance.
func NewMockTemplater(ctrl *gomock.Controller) *MockTemplater {
	mock := &MockTemplater{ctrl: ctrl}
	mock.recorder = &MockTemplaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTemplater) EXPECT() *MockTemplaterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTemplater) Create(ctx context.Context, template *models.TemplateCreate) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, template)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTemplaterMockRecorder) Create(ctx, template any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTemplater)(nil).Create), ctx, template)
}

// Delete mocks base method.
func (m *MockTemplater) Delete(ctx context.Context, id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTemplaterMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTemplater)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockTemplater) GetByID(ctx context.Context, id uint64) (*models.Template, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Template)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockTemplaterMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTemplater)(nil).GetByID), ctx, id)
}

// Update mocks base method.
func (m *MockTemplater) Update(ctx context.Context, template *models.TemplateUpdate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, template)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTemplaterMockRecorder) Update(ctx, template any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTemplater)(nil).Update), ctx, template)
}
