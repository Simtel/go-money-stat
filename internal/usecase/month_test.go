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

func (m *MockRepository) GetAll() []model.Transaction {
	args := m.Called()
	return args.Get(0).([]model.Transaction)
}

func (m *MockRepository) GetByYear(year int) []model.Transaction {
	args := m.Called(year)
	return args.Get(0).([]model.Transaction)
}

func TestGetMonthStat(t *testing.T) {
	// Create a mock repository
	mockRepo := &MockRepository{}

	// Create a Month instance
	month := NewMonth(mockRepo)

	t.Run("Get current month stat", func(t *testing.T) {
		// Arrange
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

		// Act
		monthStat := month.GetMonthStat("current")

		// Assert
		assert.Equal(t, 3, len(monthStat.Transactions))
		assert.Equal(t, 150.0, monthStat.OutComeSumm)
		assert.Equal(t, 200.0, monthStat.InComeSumm)
		assert.Equal(t, 3, monthStat.Count)
	})

	t.Run("Get previous month stat", func(t *testing.T) {
		// Arrange
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

		// Act
		monthStat := month.GetMonthStat("previous")

		// Assert
		assert.Equal(t, 3, len(monthStat.Transactions))
		assert.Equal(t, 100.0, monthStat.OutComeSumm)
		assert.Equal(t, 150.0, monthStat.InComeSumm)
		assert.Equal(t, 3, monthStat.Count)
	})
}
