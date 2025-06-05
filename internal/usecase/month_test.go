package usecase

import (
	"github.com/stretchr/testify/mock"
	"money-stat/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetCurrentMonth() []model.Transaction {
	args := m.Called()
	return args.Get(0).([]model.Transaction)
}

func (m *MockRepository) GetPreviousMonth() []model.Transaction {
	args := m.Called()
	return args.Get(0).([]model.Transaction)
}

func (m *MockRepository) GetBetweenDate(first time.Time, last time.Time) []model.Transaction {
	args := m.Called(first, last)
	return args.Get(0).([]model.Transaction)
}

func (m *MockRepository) GetAll() ([]model.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]model.Transaction), nil
}

func (m *MockRepository) GetByYear(year int) []model.Transaction {
	args := m.Called(year)
	return args.Get(0).([]model.Transaction)
}

func TestGetMonthStat(t *testing.T) {
	mockRepo := &MockRepository{}
	month := NewMonth(mockRepo)

	t.Run("Get current month stat", func(t *testing.T) {
		mockRepo.On("GetCurrentMonth").Return([]model.Transaction{
			{
				Date:       "2023-04-01",
				Outcome:    100.0,
				Income:     0.0,
				Tag:        []model.Tag{{Title: "Grocery"}},
				OutAccount: model.Account{Title: "Debit Card"},
				InAccount:  model.Account{Title: ""},
				Created:    time.Now().Unix(),
			},
			{
				Date:       "2023-04-05",
				Outcome:    50.0,
				Income:     0.0,
				Tag:        []model.Tag{{Title: "Utilities"}},
				OutAccount: model.Account{Title: "Debit Card"},
				InAccount:  model.Account{Title: ""},
				Created:    time.Now().Unix(),
			},
			{
				Date:       "2023-04-10",
				Outcome:    0.0,
				Income:     200.0,
				Tag:        []model.Tag{{Title: "Salary"}},
				OutAccount: model.Account{Title: ""},
				InAccount:  model.Account{Title: "Checking"},
				Created:    time.Now().Unix(),
			},
		})

		monthStat, _ := month.GetMonthStat("current")

		assert.Equal(t, 3, len(monthStat.Transactions))
		assert.Equal(t, 150.0, monthStat.OutcomeSumm)
		assert.Equal(t, 200.0, monthStat.IncomeSumm)
		assert.Equal(t, 3, monthStat.Count)
	})

	t.Run("Get previous month stat", func(t *testing.T) {
		mockRepo.On("GetPreviousMonth").Return([]model.Transaction{
			{
				Date:       "2023-03-01",
				Outcome:    75.0,
				Income:     0.0,
				Tag:        []model.Tag{{Title: "Rent"}},
				OutAccount: model.Account{Title: "Debit Card"},
				InAccount:  model.Account{Title: ""},
				Created:    time.Now().Unix(),
			},
			{
				Date:       "2023-03-15",
				Outcome:    25.0,
				Income:     0.0,
				Tag:        []model.Tag{{Title: "Groceries"}},
				OutAccount: model.Account{Title: "Debit Card"},
				InAccount:  model.Account{Title: ""},
				Created:    time.Now().Unix(),
			},
			{
				Date:       "2023-03-20",
				Outcome:    0.0,
				Income:     150.0,
				Tag:        []model.Tag{{Title: "Freelance"}},
				OutAccount: model.Account{Title: ""},
				InAccount:  model.Account{Title: "Checking"},
				Created:    time.Now().Unix(),
			},
		})

		monthStat, _ := month.GetMonthStat("previous")

		assert.Equal(t, 3, len(monthStat.Transactions))
		assert.Equal(t, 100.0, monthStat.OutcomeSumm)
		assert.Equal(t, 150.0, monthStat.IncomeSumm)
		assert.Equal(t, 3, monthStat.Count)
	})
}

func getMonth() *Month {

	return NewMonth(&MockRepository{})
}

func TestGetAccountTitle(t *testing.T) {
	month := getMonth()

	accountIn := model.Account{Title: "Debit Card"}
	accountOut := model.Account{Title: "Credit Card"}

	transactionIn := model.Transaction{InAccount: accountIn, Income: 100}
	transactionOut := model.Transaction{OutAccount: accountOut, Outcome: 100}

	assert.Equal(t, "Debit Card", month.getAccountTitle(transactionIn))
	assert.Equal(t, "Credit Card", month.getAccountTitle(transactionOut))
}
