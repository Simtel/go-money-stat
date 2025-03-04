package app

import (
	"errors"
	"testing"
)

func TestGetGlobalApp(t *testing.T) {

	testCases := []struct {
		name     string
		app      *App
		expected error
	}{
		{
			name:     "Global app is not initialized",
			app:      nil,
			expected: errors.New("global app is not initialized"),
		},
		{
			name:     "Global app is initialized",
			app:      &App{},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetGlobalApp(tc.app)

			_, err := GetGlobalApp()
			if err != nil {
				if err.Error() != tc.expected.Error() {
					t.Errorf("Ожидалась ошибка: %v, но получили: %v", tc.expected, err)
				}
			}
		})
	}
}
