package usecase_test

import (
	"money-stat/internal/adapter/db"
	"money-stat/internal/model"
	"money-stat/internal/services/zenmoney"
	"money-stat/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) (tx db.DBServiceInterface) {
	m.Called(value)
	return m
}

func (m *MockDB) Where(query interface{}, args ...interface{}) (tx db.DBServiceInterface) {
	m.Called(query, args)
	return m
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) (tx db.DBServiceInterface) {
	if len(conds) > 0 {
		m.Called(value, conds)
	}
	m.Called(value)
	return m
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) db.DBServiceInterface {
	m.Called(dest, conds)
	return m
}

func (m *MockDB) Updates(value interface{}) db.DBServiceInterface {
	m.Called(value)
	return m
}

func (m *MockDB) Exec(sql string, values ...interface{}) db.DBServiceInterface {
	m.Called(sql, values)
	return m
}

func (m *MockDB) GetDB() *gorm.DB {
	// Возвращаем пустую, но валидную DB для тестов
	db := &gorm.DB{}
	return db
}

func (m *MockDB) Model(dest interface{}) db.DBServiceInterface {
	m.Called(dest)
	return m
}

func (m *MockDB) Association(field string) *gorm.Association {
	// Возвращаем nil для тестов (не критично для тестов)
	return nil
}

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) Diff() (*zenmoney.Response, error) {
	args := m.Called()
	return args.Get(0).(*zenmoney.Response), args.Error(1)
}

func (m *MockAPI) DiffSince(timestamp int64) (*zenmoney.Response, error) {
	args := m.Called(timestamp)
	return args.Get(0).(*zenmoney.Response), args.Error(1)
}

// Тест только для создания Sync - полная логика тестируется интеграционно
func TestNewSync(t *testing.T) {
	mockDB := new(MockDB)
	mockAPI := &zenmoney.Api{}

	sync := usecase.NewSync(mockDB, mockAPI)

	assert.NotNil(t, sync, "NewSync should return a non-nil Sync instance")
}

func TestSync_ClearTables(t *testing.T) {
	mockDB := new(MockDB)
	mockAPI := new(MockAPI)

	mockDB.On("Where", "`id` != ?", []interface{}{""}).Return(mockDB)
	mockDB.On("Delete", &model.Transaction{}).Return(mockDB)
	mockDB.On("Delete", &model.Account{}).Return(mockDB)
	mockDB.On("Delete", &model.Tag{}).Return(mockDB)
	mockDB.On("Delete", &model.Instrument{}).Return(mockDB)

	sync := usecase.NewSync(mockDB, mockAPI)

	sync.ClearTables()

	mockAPI.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}
