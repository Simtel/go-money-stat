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
			// Очищаем окружение перед тестом
			_ = os.Unsetenv("ZENMONEY_TOKEN")

			for key, value := range tc.env {
				err := os.Setenv(key, value)
				if err != nil {
					t.Fatalf("Failed to set env: %v", err)
				}
			}

			// Очищаем окружение после теста
			t.Cleanup(func() {
				for key := range tc.env {
					_ = os.Unsetenv(key)
				}
			})

			config := New()

			if config.ZenMoney.Token != tc.expected {
				t.Errorf("Ожидалось, что token будет %v, но получили %v", tc.expected, config.ZenMoney.Token)
			}
		})
	}
}
