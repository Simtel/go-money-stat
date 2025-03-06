package mock_transactions

import (
	"github.com/golang/mock/gomock"
	"money-stat/internal/model"
	"reflect"
	"time"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetBetweenDate mocks base method.
func (m *MockRepositoryInterface) GetBetweenDate(first, last time.Time) []model.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBetweenDate", first, last)
	ret0, _ := ret[0].([]model.Transaction)
	return ret0
}

// GetBetweenDate indicates an expected call of GetBetweenDate.
func (mr *MockRepositoryInterfaceMockRecorder) GetBetweenDate(first, last interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBetweenDate", reflect.TypeOf((*MockRepositoryInterface)(nil).GetBetweenDate), first, last)
}

// GetCurrentMonth mocks base method.
func (m *MockRepositoryInterface) GetCurrentMonth() []model.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentMonth")
	ret0, _ := ret[0].([]model.Transaction)
	return ret0
}

// GetCurrentMonth indicates an expected call of GetCurrentMonth.
func (mr *MockRepositoryInterfaceMockRecorder) GetCurrentMonth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentMonth", reflect.TypeOf((*MockRepositoryInterface)(nil).GetCurrentMonth))
}

// GetPreviousMonth mocks base method.
func (m *MockRepositoryInterface) GetPreviousMonth() []model.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreviousMonth")
	ret0, _ := ret[0].([]model.Transaction)
	return ret0
}

// GetPreviousMonth indicates an expected call of GetPreviousMonth.
func (mr *MockRepositoryInterfaceMockRecorder) GetPreviousMonth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreviousMonth", reflect.TypeOf((*MockRepositoryInterface)(nil).GetPreviousMonth))
}
