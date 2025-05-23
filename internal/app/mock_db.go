package app

import (
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
	"reflect"
)

type MockDbInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDbInterfaceMockRecorder
}

// MockDbInterfaceMockRecorder is the mock recorder for MockDbInterface.
type MockDbInterfaceMockRecorder struct {
	mock *MockDbInterface
}

// NewMockDbInterface creates a new mock instance.
func NewMockDbInterface(ctrl *gomock.Controller) *MockDbInterface {
	mock := &MockDbInterface{ctrl: ctrl}
	mock.recorder = &MockDbInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDbInterface) EXPECT() *MockDbInterfaceMockRecorder {
	return m.recorder
}

// GetGorm mocks base method.
func (m *MockDbInterface) GetGorm() *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGorm")
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// GetGorm indicates an expected call of GetGorm.
func (mr *MockDbInterfaceMockRecorder) GetGorm() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGorm", reflect.TypeOf((*MockDbInterface)(nil).GetGorm))
}
