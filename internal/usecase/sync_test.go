package usecase_test

import (
	"money-stat/internal/adapter/db"
	"money-stat/internal/model"
	"money-stat/internal/services/zenmoney"
	"money-stat/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) (tx db.DBServiceInterface) {
	//m.Called(value)
	return m
}

func (m *MockDB) Where(query interface{}, args ...interface{}) (tx db.DBServiceInterface) {
	//m.Called(query, args)
	return m
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) (tx db.DBServiceInterface) {
	///m.Called(value, conds)
	return m
}

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) Diff() (*zenmoney.Response, error) {
	args := m.Called()
	return args.Get(0).(*zenmoney.Response), args.Error(1)
}

type MockSpinner struct {
	mock.Mock
}

func (m *MockSpinner) Success(text string) {
	m.Called(text)
}

func TestSync_FullSync(t *testing.T) {
	mockDB := new(MockDB)
	mockAPI := new(MockAPI)

	diffResponse := &zenmoney.Response{
		Tag: []zenmoney.Tag{
			{Id: "tag1", Title: "Tag 1"},
			{Id: "tag2", Title: "Tag 2"},
		},
		Instrument: []zenmoney.Instrument{
			{Id: 1, Title: "Instrument 1", ShortTitle: "I1", Symbol: "$", Rate: 1.0},
		},
		Account: []zenmoney.Account{
			{Id: "acc1", Title: "Account 1", Balance: 1000, Instrument: 1, StartBalance: 0},
		},
		Transaction: []zenmoney.Transaction{
			{
				Id:                "trans1",
				Changed:           time.Now().Unix(),
				Created:           time.Now().Unix(),
				IncomeInstrument:  1,
				Income:            100,
				OutcomeInstrument: 1,
				Outcome:           100,
				Date:              time.Now().Format("2006-01-02"),
				Deleted:           false,
				IncomeAccount:     "acc1",
				OutcomeAccount:    "acc1",
				Tag:               []string{"tag1"},
				Comment:           "Test transaction",
			},
		},
	}

	mockAPI.On("Diff").Return(diffResponse, nil)

	mockDB.On("Where", "`id` != ?", "").Return(mockDB)
	mockDB.On("Delete", &model.Transaction{}).Return(&gorm.DB{})
	mockDB.On("Delete", &model.Account{}).Return(&gorm.DB{})
	mockDB.On("Delete", &model.Tag{}).Return(&gorm.DB{})
	mockDB.On("Delete", &model.Instrument{}).Return(&gorm.DB{})

	mockDB.On("Create", mock.AnythingOfType("*model.Tag")).Return(&gorm.DB{})
	mockDB.On("Create", mock.AnythingOfType("*model.Instrument")).Return(&gorm.DB{})
	mockDB.On("Create", mock.AnythingOfType("*model.Account")).Return(&gorm.DB{})
	mockDB.On("Create", mock.AnythingOfType("*model.Transaction")).Return(&gorm.DB{})

	sync := usecase.NewSync(mockDB, mockAPI)

	sync.FullSync()

}

func TestSync_ClearTables(t *testing.T) {

	mockDB := new(MockDB)
	mockAPI := new(MockAPI)

	mockDB.On("Where", "`id` != ?", "").Return(mockDB)
	mockDB.On("Delete", &model.Transaction{}).Return(&gorm.DB{})
	mockDB.On("Delete", &model.Account{}).Return(&gorm.DB{})
	mockDB.On("Delete", &model.Tag{}).Return(&gorm.DB{})
	mockDB.On("Delete", &model.Instrument{}).Return(&gorm.DB{})

	sync := usecase.NewSync(mockDB, mockAPI)

	sync.ClearTables()

}

func TestNewSync(t *testing.T) {
	mockDB := new(MockDB)
	mockAPI := &zenmoney.Api{}

	sync := usecase.NewSync(mockDB, mockAPI)

	assert.NotNil(t, sync, "NewSync should return a non-nil Sync instance")
}
