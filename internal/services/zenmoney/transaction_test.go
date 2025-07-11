package zenmoney

import (
	"testing"
)

func TestTransaction_FormatAmount(t *testing.T) {
	tests := []struct {
		name        string
		transaction Transaction
		expected    string
	}{
		{
			name: "Only outcome",
			transaction: Transaction{
				Income:  0,
				Outcome: 150.33,
			},
			expected: "-150.33",
		},
		{
			name: "Only income",
			transaction: Transaction{
				Income:  200.44,
				Outcome: 0,
			},
			expected: "200.44",
		},
		{
			name: "Transfer",
			transaction: Transaction{
				Income:  300.55,
				Outcome: 100.22,
			},
			expected: "100.22 -> 300.55",
		},
		{
			name: "No amount",
			transaction: Transaction{
				Income:  0,
				Outcome: 0,
			},
			expected: "0",
		},
		{
			name: "Zero income with zero outcome",
			transaction: Transaction{
				Income:  0,
				Outcome: 0,
			},
			expected: "0",
		},
		{
			name: "Zero income with positive outcome",
			transaction: Transaction{
				Income:  0,
				Outcome: 50.75,
			},
			expected: "-50.75",
		},
		{
			name: "Positive income with zero outcome",
			transaction: Transaction{
				Income:  100.0,
				Outcome: 0,
			},
			expected: "100.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.transaction.FormatAmount()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestTransaction_IsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		deleted  bool
		expected bool
	}{
		{name: "Not deleted", deleted: false, expected: false},
		{name: "Deleted", deleted: true, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := Transaction{Deleted: tt.deleted}
			if result := transaction.IsDeleted(); result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTransaction_IsIncome(t *testing.T) {
	tests := []struct {
		name     string
		income   float64
		outcome  float64
		expected bool
	}{
		{name: "Positive income and zero outcome", income: 100.0, outcome: 0, expected: true},
		{name: "Zero income and zero outcome", income: 0, outcome: 0, expected: false},
		{name: "Positive income and positive outcome", income: 100.0, outcome: 50.0, expected: false},
		{name: "Zero income and positive outcome", income: 0, outcome: 50.0, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := Transaction{Income: tt.income, Outcome: tt.outcome}
			if result := transaction.IsIncome(); result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTransaction_IsOutcome(t *testing.T) {
	tests := []struct {
		name     string
		income   float64
		outcome  float64
		expected bool
	}{
		{name: "Positive outcome and zero income", income: 0, outcome: 100.0, expected: true},
		{name: "Zero income and zero outcome", income: 0, outcome: 0, expected: false},
		{name: "Positive income and positive outcome", income: 100.0, outcome: 50.0, expected: false},
		{name: "Positive income and zero outcome", income: 100.0, outcome: 0, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := Transaction{Income: tt.income, Outcome: tt.outcome}
			if result := transaction.IsOutcome(); result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTransaction_IsTransfer(t *testing.T) {
	tests := []struct {
		name     string
		income   float64
		outcome  float64
		expected bool
	}{
		{name: "Positive income and positive outcome", income: 100.0, outcome: 50.0, expected: true},
		{name: "Zero income and zero outcome", income: 0, outcome: 0, expected: false},
		{name: "Positive income and zero outcome", income: 100.0, outcome: 0, expected: false},
		{name: "Zero income and positive outcome", income: 0, outcome: 100.0, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := Transaction{Income: tt.income, Outcome: tt.outcome}
			if result := transaction.IsTransfer(); result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
