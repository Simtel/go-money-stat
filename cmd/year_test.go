package cmd

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestRunYear(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := RunYear()

	// Тестирование пограничных случаев
	testCases := []struct {
		name     string
		args     []string
		expected error
	}{
		{
			name:     "Неверный год",
			args:     []string{"2019"},
			expected: fmt.Errorf("Указан неверный год: 2019"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd.SetArgs(tc.args)

			err := cmd.Execute()

			if err.Error() != tc.expected.Error() {
				t.Errorf("Ожидалась ошибка: %v, но получили: %v", tc.expected, err)
			}
		})
	}
}
