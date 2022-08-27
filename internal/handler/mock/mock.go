// Code generated by MockGen. DO NOT EDIT.
// Source: internal/handler/interfaces.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/nnemirovsky/pacgen/internal/model"
)

// ProxyProfileService is a mock of ProxyProfileService interface.
type ProxyProfileService struct {
	ctrl     *gomock.Controller
	recorder *ProxyProfileServiceMockRecorder
}

// ProxyProfileServiceMockRecorder is the mock recorder for ProxyProfileService.
type ProxyProfileServiceMockRecorder struct {
	mock *ProxyProfileService
}

// NewProxyProfileService creates a new mock instance.
func NewProxyProfileService(ctrl *gomock.Controller) *ProxyProfileService {
	mock := &ProxyProfileService{ctrl: ctrl}
	mock.recorder = &ProxyProfileServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *ProxyProfileService) EXPECT() *ProxyProfileServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *ProxyProfileService) Create(ctx context.Context, profile *model.ProxyProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *ProxyProfileServiceMockRecorder) Create(ctx, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*ProxyProfileService)(nil).Create), ctx, profile)
}

// Delete mocks base method.
func (m *ProxyProfileService) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *ProxyProfileServiceMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*ProxyProfileService)(nil).Delete), ctx, id)
}

// GetAll mocks base method.
func (m *ProxyProfileService) GetAll(ctx context.Context) ([]model.ProxyProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]model.ProxyProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *ProxyProfileServiceMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*ProxyProfileService)(nil).GetAll), ctx)
}

// GetByID mocks base method.
func (m *ProxyProfileService) GetByID(ctx context.Context, id int) (model.ProxyProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(model.ProxyProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *ProxyProfileServiceMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*ProxyProfileService)(nil).GetByID), ctx, id)
}

// Update mocks base method.
func (m *ProxyProfileService) Update(ctx context.Context, profile model.ProxyProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *ProxyProfileServiceMockRecorder) Update(ctx, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*ProxyProfileService)(nil).Update), ctx, profile)
}

// RuleService is a mock of RuleService interface.
type RuleService struct {
	ctrl     *gomock.Controller
	recorder *RuleServiceMockRecorder
}

// RuleServiceMockRecorder is the mock recorder for RuleService.
type RuleServiceMockRecorder struct {
	mock *RuleService
}

// NewRuleService creates a new mock instance.
func NewRuleService(ctrl *gomock.Controller) *RuleService {
	mock := &RuleService{ctrl: ctrl}
	mock.recorder = &RuleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *RuleService) EXPECT() *RuleServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *RuleService) Create(ctx context.Context, rule *model.Rule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, rule)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *RuleServiceMockRecorder) Create(ctx, rule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*RuleService)(nil).Create), ctx, rule)
}

// Delete mocks base method.
func (m *RuleService) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *RuleServiceMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*RuleService)(nil).Delete), ctx, id)
}

// GetAll mocks base method.
func (m *RuleService) GetAll(ctx context.Context) ([]model.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]model.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *RuleServiceMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*RuleService)(nil).GetAll), ctx)
}

// GetAllWithProfiles mocks base method.
func (m *RuleService) GetAllWithProfiles(ctx context.Context) ([]model.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllWithProfiles", ctx)
	ret0, _ := ret[0].([]model.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllWithProfiles indicates an expected call of GetAllWithProfiles.
func (mr *RuleServiceMockRecorder) GetAllWithProfiles(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllWithProfiles", reflect.TypeOf((*RuleService)(nil).GetAllWithProfiles), ctx)
}

// GetByID mocks base method.
func (m *RuleService) GetByID(ctx context.Context, id int) (model.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(model.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *RuleServiceMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*RuleService)(nil).GetByID), ctx, id)
}

// Update mocks base method.
func (m *RuleService) Update(ctx context.Context, rule model.Rule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, rule)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *RuleServiceMockRecorder) Update(ctx, rule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*RuleService)(nil).Update), ctx, rule)
}
