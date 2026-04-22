package cmd

import (
	"fmt"
	"money-stat/internal/app"
	"testing"
)

func TestRunYear(t *testing.T) {
	app := app.NewApp(app.NewContainer(app.NewDB()))

	cmd := RunYear(app)

	// Тестирование пограничных случаев
	testCases := []struct {
		name     string
		args     []string
		expected error
	}{
		{
			name:     "Неверный год",
			args:     []string{"2019"},
			expected: fmt.Errorf("указан неверный год: 2019"),
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
