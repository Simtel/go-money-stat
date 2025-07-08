package usecase_test

import (
	"github.com/stretchr/testify/mock"
	"money-stat/internal/model"
	"money-stat/internal/usecase"
	"reflect"
	"testing"
	"time"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) GetCurrentMonth() []model.Transaction {
	args := m.Called()
	return args.Get(0).([]model.Transaction)
}

func (m *mockRepository) GetPreviousMonth() []model.Transaction {
	args := m.Called()
	return args.Get(0).([]model.Transaction)
}

func (m *mockRepository) GetBetweenDate(first time.Time, last time.Time) []model.Transaction {
	args := m.Called(first, last)
	return args.Get(0).([]model.Transaction)
}

func (m *mockRepository) GetAll() ([]model.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]model.Transaction), nil
}

func (m *mockRepository) GetByYear(year int) []model.Transaction {
	args := m.Called(year)
	return args.Get(0).([]model.Transaction)
}

func TestYear_GetYearStat(t *testing.T) {
	testCases := []struct {
		name          string
		transactions  []model.Transaction
		selectYear    int
		expectedStats []usecase.MonthStat
	}{
		{
			name: "Single transaction",
			transactions: []model.Transaction{
				{Date: "2023-03-01", Outcome: 100.0, Income: 0.0},
			},
			selectYear: 2023,
			expectedStats: []usecase.MonthStat{
				{Month: "2023-03", Income: 0.0, OutCome: 100.0},
			},
		},
		{
			name: "Multiple transactions",
			transactions: []model.Transaction{
				{Date: "2023-03-01", Outcome: 100.0, Income: 0.0},
				{Date: "2023-03-15", Income: 200.0, Outcome: 0.0},
				{Date: "2023-04-01", Outcome: 50.0, Income: 0.0},
				{Date: "2023-04-15", Income: 300.0, Outcome: 0.0},
			},
			selectYear: 2023,
			expectedStats: []usecase.MonthStat{
				{Month: "2023-03", Income: 200.0, OutCome: 100.0},
				{Month: "2023-04", Income: 300.0, OutCome: 50.0},
			},
		},
		{
			name:          "No transactions",
			transactions:  []model.Transaction{},
			selectYear:    2023,
			expectedStats: []usecase.MonthStat{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockRepository{}
			year := usecase.NewYear(repo)
			stats := year.GetYearStat(tc.selectYear)

			if !reflect.DeepEqual(stats, tc.expectedStats) {
				t.Errorf("GetYearStat() = %v, expected %v", stats, tc.expectedStats)
			}
		})
	}
}
