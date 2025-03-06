package usecase

import (
	"errors"
	"github.com/golang/mock/gomock"
	"money-stat/internal/adapter/sqliterepo/zenrepo/transactions/mock_transactions"
	"money-stat/internal/model"
	"testing"
)

func TestMonth_GetMonthStat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_transactions.NewMockRepositoryInterface(ctrl)

	month := NewMonth(mockRepo)

	// Тестирование пограничных случаев
	testCases := []struct {
		name     string
		month    string
		expected error
	}{
		{
			name:     "Current month",
			month:    "current",
			expected: nil,
		},
		{
			name:     "Previous month",
			month:    "previous",
			expected: nil,
		},
		{
			name:     "Invalid month",
			month:    "invalid",
			expected: errors.New("invalid month"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Устанавливаем ожидаемые результаты для мока
			if tc.month == "current" {
				mockRepo.EXPECT().GetCurrentMonth().Return([]model.Transaction{})
			}

			if tc.month == "previous" {
				mockRepo.EXPECT().GetPreviousMonth().Return([]model.Transaction{})
			}

			// Вызываем метод GetMonthStat
			month.GetMonthStat(tc.month)

		})
	}
}
