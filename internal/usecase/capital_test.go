package usecase

import (
	"errors"
	"money-stat/internal/model"
	"reflect"
	"testing"
	"time"
)

type MockAccountRepository struct {
	accounts []model.Account
}

func (m *MockAccountRepository) GetAll() []model.Account {
	return m.accounts
}

type MockTransactionRepository struct {
	transactions []model.Transaction
	err          error
}

func (m *MockTransactionRepository) GetAll() ([]model.Transaction, error) {
	return m.transactions, m.err
}

func (m *MockTransactionRepository) GetCurrentMonth() []model.Transaction {
	return m.transactions
}

func (m *MockTransactionRepository) GetPreviousMonth() []model.Transaction {
	return m.transactions
}

func (m *MockTransactionRepository) GetBetweenDate(first time.Time, last time.Time) []model.Transaction {
	return m.transactions
}

func (m *MockTransactionRepository) GetByYear(year int) []model.Transaction {
	return m.transactions
}

func TestNewCapital(t *testing.T) {
	mockTransRepo := &MockTransactionRepository{}
	mockAccRepo := &MockAccountRepository{}

	capital := NewCapital(mockTransRepo, mockAccRepo)

	if capital.transactionRepo != mockTransRepo {
		t.Error("Transaction repository not properly initialized")
	}

	if capital.accountRepo != mockAccRepo {
		t.Error("Account repository not properly initialized")
	}
}

func TestCapital_GetCapital_EmptyData(t *testing.T) {
	mockTransRepo := &MockTransactionRepository{
		transactions: []model.Transaction{},
	}
	mockAccRepo := &MockAccountRepository{
		accounts: []model.Account{},
	}

	capital := NewCapital(mockTransRepo, mockAccRepo)
	result, err := capital.GetCapital()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestCapital_GetCapital_WithTransactions(t *testing.T) {
	today := time.Now().Format(dateLayout)
	nextMonth := time.Now().AddDate(0, 1, 0).Format(dateLayout)

	todayMonth, _ := time.Parse(dateLayout, today)
	todayMonthKey := todayMonth.Format(monthKeyLayout)

	nextMonthParsed, _ := time.Parse(dateLayout, nextMonth)
	nextMonthKey := nextMonthParsed.Format(monthKeyLayout)

	mockTransRepo := &MockTransactionRepository{
		transactions: []model.Transaction{
			{
				Date:   today,
				Income: 500,
				InAccount: model.Account{
					Currency: model.Instrument{ShortTitle: "RUB", Rate: 1.0},
				},
			},
			{
				Date:    nextMonth,
				Outcome: 200,
				OutAccount: model.Account{
					Currency: model.Instrument{ShortTitle: "RUB", Rate: 1.0},
				},
			},
		},
	}
	mockAccRepo := &MockAccountRepository{
		accounts: []model.Account{
			{
				Id:           "1",
				StartBalance: 1000,
				Currency: model.Instrument{
					ShortTitle: "RUB",
					Rate:       1.0,
				},
			},
		},
	}

	capital := NewCapital(mockTransRepo, mockAccRepo)
	result, err := capital.GetCapital()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []CapitalDto{
		{
			Month:   todayMonthKey,
			Balance: 500,
		},
		{
			Month:   nextMonthKey,
			Balance: -200,
		},
	}

	if len(result) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(result))
	}

	monthsFound := make(map[string]bool)
	for _, item := range result {
		monthsFound[item.Month] = true
	}

	if !monthsFound[todayMonthKey] {
		t.Errorf("Current month %s not found in results", todayMonthKey)
	}
	if !monthsFound[nextMonthKey] {
		t.Errorf("Next month %s not found in results", nextMonthKey)
	}
}

func TestCapital_GetCapital_TransactionRepoError(t *testing.T) {
	expectedErr := errors.New("database error")
	mockTransRepo := &MockTransactionRepository{
		transactions: nil,
		err:          expectedErr,
	}
	mockAccRepo := &MockAccountRepository{
		accounts: []model.Account{},
	}

	capital := NewCapital(mockTransRepo, mockAccRepo)
	_, err := capital.GetCapital()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "failed to process transactions: failed to fetch transactions: database error" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCapital_ExtractMonthKey(t *testing.T) {
	capital := &Capital{}

	tests := []struct {
		name     string
		dateStr  string
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid date",
			dateStr:  "2023-05-15",
			expected: "2023-05",
			wantErr:  false,
		},
		{
			name:     "Invalid date format",
			dateStr:  "15/05/2023",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := capital.extractMonthKey(tt.dateStr)

			if (err != nil) != tt.wantErr {
				t.Errorf("extractMonthKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected {
				t.Errorf("extractMonthKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCapital_GetOrCreateMonthlyStat(t *testing.T) {
	capital := &Capital{}
	monthlyStats := map[string]CapitalDto{
		"2023-01": {Month: "2023-01", Balance: 1000},
	}

	tests := []struct {
		name     string
		monthKey string
		expected CapitalDto
	}{
		{
			name:     "Existing month",
			monthKey: "2023-01",
			expected: CapitalDto{Month: "2023-01", Balance: 1000},
		},
		{
			name:     "New month",
			monthKey: "2023-02",
			expected: CapitalDto{Month: "2023-02", Balance: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := capital.getOrCreateMonthlyStat(monthlyStats, tt.monthKey)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getOrCreateMonthlyStat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCapital_ApplyTransactionToBalance(t *testing.T) {
	capital := &Capital{}
	stat := CapitalDto{Month: "2023-01", Balance: 1000}

	tests := []struct {
		name        string
		transaction model.Transaction
		expected    CapitalDto
	}{
		{
			name: "Income transaction",
			transaction: model.Transaction{
				Income: 500,
				InAccount: model.Account{
					Currency: model.Instrument{ShortTitle: "RUB", Rate: 1.0},
				},
			},
			expected: CapitalDto{Month: "2023-01", Balance: 1500},
		},
		{
			name: "Outcome transaction",
			transaction: model.Transaction{
				Outcome: 300,
				OutAccount: model.Account{
					Currency: model.Instrument{ShortTitle: "RUB", Rate: 1.0},
				},
			},
			expected: CapitalDto{Month: "2023-01", Balance: 700},
		},
		{
			name: "Foreign currency income",
			transaction: model.Transaction{
				Income: 100,
				InAccount: model.Account{
					Currency: model.Instrument{ShortTitle: "USD", Rate: 75.0},
				},
			},
			expected: CapitalDto{Month: "2023-01", Balance: 8500}, // 1000 + (100 * 75)
		},
		{
			name: "Foreign currency outcome",
			transaction: model.Transaction{
				Outcome: 10,
				OutAccount: model.Account{
					Currency: model.Instrument{ShortTitle: "USD", Rate: 75.0},
				},
			},
			expected: CapitalDto{Month: "2023-01", Balance: 250}, // 1000 - (10 * 75)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := capital.applyTransactionToBalance(stat, tt.transaction)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("applyTransactionToBalance() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCapital_ConvertToRubles(t *testing.T) {
	capital := &Capital{}

	tests := []struct {
		name     string
		amount   float64
		account  model.Account
		expected float64
	}{
		{
			name:   "Ruble account",
			amount: 1000,
			account: model.Account{
				Currency: model.Instrument{ShortTitle: "RUB", Rate: 1.0},
			},
			expected: 1000,
		},
		{
			name:   "Dollar account",
			amount: 100,
			account: model.Account{
				Currency: model.Instrument{ShortTitle: "USD", Rate: 75.0},
			},
			expected: 7500,
		},
		{
			name:   "Euro account",
			amount: 50,
			account: model.Account{
				Currency: model.Instrument{ShortTitle: "EUR", Rate: 85.0},
			},
			expected: 4250,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := capital.convertToRubles(tt.amount, tt.account)

			if result != tt.expected {
				t.Errorf("convertToRubles() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCapital_ConvertToSortedSlice(t *testing.T) {
	capital := &Capital{}
	monthlyStats := map[string]CapitalDto{
		"2023-03": {Month: "2023-03", Balance: 3000},
		"2023-01": {Month: "2023-01", Balance: 1000},
		"2023-02": {Month: "2023-02", Balance: 2000},
	}

	expected := []CapitalDto{
		{Month: "2023-01", Balance: 1000},
		{Month: "2023-02", Balance: 2000},
		{Month: "2023-03", Balance: 3000},
	}

	result := capital.convertToSortedSlice(monthlyStats)

	if len(result) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(result))
		return
	}

	for i, item := range result {
		if item.Month != expected[i].Month || item.Balance != expected[i].Balance {
			t.Errorf("Item %d: expected %v, got %v", i, expected[i], item)
		}
	}
}

func TestCapital_IsOutcome(t *testing.T) {
	tests := []struct {
		name        string
		transaction model.Transaction
		expected    bool
	}{
		{
			name: "Is outcome",
			transaction: model.Transaction{
				Outcome: 100,
			},
			expected: true,
		},
		{
			name: "Not outcome",
			transaction: model.Transaction{
				Outcome: 0,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.transaction.IsOutcome()

			if result != tt.expected {
				t.Errorf("IsOutcome() = %v, want %v", result, tt.expected)
			}
		})
	}
}
