package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Тестирование пограничных случаев
	testCases := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "Environment variable is not set",
			env:      map[string]string{},
			expected: "",
		},
		{
			name:     "Environment variable is set",
			env:      map[string]string{"ZENMONEY_TOKEN": "test_token"},
			expected: "test_token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.env {
				os.Setenv(key, value)
			}

			config := New()

			if config.ZenMoney.Token != tc.expected {
				t.Errorf("Ожидалось, что token будет %v, но получили %v", tc.expected, config.ZenMoney.Token)
			}
		})
	}
}
