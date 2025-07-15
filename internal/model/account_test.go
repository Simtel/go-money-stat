package model

import (
	"testing"
)

func TestAccount_IsRuble(t *testing.T) {
	tests := []struct {
		name       string
		instrument int
		expected   bool
	}{
		{
			name:       "Ruble account",
			instrument: 2,
			expected:   true,
		},
		{
			name:       "Dollar account",
			instrument: 1,
			expected:   false,
		},
		{
			name:       "Other currency",
			instrument: 3,
			expected:   false,
		},
		{
			name:       "Zero instrument",
			instrument: 0,
			expected:   false,
		},
		{
			name:       "Negative instrument",
			instrument: -1,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Instrument: tt.instrument,
			}
			result := account.IsRuble()
			if result != tt.expected {
				t.Errorf("IsRuble() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestAccount_IsDollar(t *testing.T) {
	tests := []struct {
		name       string
		instrument int
		expected   bool
	}{
		{
			name:       "Dollar account",
			instrument: 1,
			expected:   true,
		},
		{
			name:       "Ruble account",
			instrument: 2,
			expected:   false,
		},
		{
			name:       "Other currency",
			instrument: 3,
			expected:   false,
		},
		{
			name:       "Zero instrument",
			instrument: 0,
			expected:   false,
		},
		{
			name:       "Negative instrument",
			instrument: -1,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Instrument: tt.instrument,
			}
			result := account.IsDollar()
			if result != tt.expected {
				t.Errorf("IsDollar() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestAccount_BothMethods(t *testing.T) {
	account := &Account{
		Id:           "test-id",
		Title:        "Test Account",
		Balance:      1000.0,
		StartBalance: 500.0,
		Instrument:   2,
	}

	if !account.IsRuble() {
		t.Error("Account with instrument 2 should be ruble")
	}
	if account.IsDollar() {
		t.Error("Account with instrument 2 should not be dollar")
	}

	account.Instrument = 1
	if account.IsRuble() {
		t.Error("Account with instrument 1 should not be ruble")
	}
	if !account.IsDollar() {
		t.Error("Account with instrument 1 should be dollar")
	}

	account.Instrument = 3
	if account.IsRuble() {
		t.Error("Account with instrument 3 should not be ruble")
	}
	if account.IsDollar() {
		t.Error("Account with instrument 3 should not be dollar")
	}
}

func TestAccount_Structure(t *testing.T) {
	account := &Account{
		Id:           "account-123",
		Title:        "My Test Account",
		Balance:      2500.50,
		StartBalance: 1000.00,
		Instrument:   1,
	}

	if account.Id != "account-123" {
		t.Errorf("Expected Id to be 'account-123', got %s", account.Id)
	}
	if account.Title != "My Test Account" {
		t.Errorf("Expected Title to be 'My Test Account', got %s", account.Title)
	}
	if account.Balance != 2500.50 {
		t.Errorf("Expected Balance to be 2500.50, got %f", account.Balance)
	}
	if account.StartBalance != 1000.00 {
		t.Errorf("Expected StartBalance to be 1000.00, got %f", account.StartBalance)
	}
	if account.Instrument != 1 {
		t.Errorf("Expected Instrument to be 1, got %d", account.Instrument)
	}
}

func TestAccount_ZeroValues(t *testing.T) {
	account := &Account{}

	if account.IsRuble() {
		t.Error("Empty account should not be ruble")
	}
	if account.IsDollar() {
		t.Error("Empty account should not be dollar")
	}
	if account.Id != "" {
		t.Error("Empty account should have empty Id")
	}
	if account.Balance != 0.0 {
		t.Error("Empty account should have zero Balance")
	}
}

func TestAccount_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		instrument int
		isRuble    bool
		isDollar   bool
	}{
		{"Max int", 2147483647, false, false},
		{"Min int", -2147483648, false, false},
		{"Just before ruble", 1, false, true},
		{"Just after ruble", 3, false, false},
		{"Just before dollar", 0, false, false},
		{"Just after dollar", 2, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{Instrument: tt.instrument}

			if account.IsRuble() != tt.isRuble {
				t.Errorf("IsRuble() = %v, expected %v for instrument %d",
					account.IsRuble(), tt.isRuble, tt.instrument)
			}

			if account.IsDollar() != tt.isDollar {
				t.Errorf("IsDollar() = %v, expected %v for instrument %d",
					account.IsDollar(), tt.isDollar, tt.instrument)
			}
		})
	}
}

func BenchmarkAccount_IsRuble(b *testing.B) {
	account := &Account{Instrument: 2}
	for i := 0; i < b.N; i++ {
		account.IsRuble()
	}
}

func BenchmarkAccount_IsDollar(b *testing.B) {
	account := &Account{Instrument: 1}
	for i := 0; i < b.N; i++ {
		account.IsDollar()
	}
}
