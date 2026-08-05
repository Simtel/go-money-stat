package menu

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"money-stat/internal/app"
)

// Start запускает интерактивное меню выбора команды
func Start(rootCmd *cobra.Command) error {
	pterm.DefaultSection.Println("Выберите команду money-stat")

	options := selectableCommands(rootCmd)
	if len(options) == 0 {
		return fmt.Errorf("нет доступных команд")
	}

	selected, err := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		WithDefaultText("Выберите команду").
		WithMaxHeight(8).
		Show()
	if err != nil {
		return err
	}

	execArgs, err := buildArgs(rootCmd, selected)
	if err != nil {
		return err
	}

	rootCmd.SetArgs(execArgs)
	return rootCmd.ExecuteContext(rootCmd.Context())
}

// Run возвращает cobra-команду интерактивного меню
func Run(_ *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "menu",
		Short: "Интерактивный выбор команды из списка",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Start(cmd.Root())
		},
	}
}

// selectableCommands возвращает имена доступных для выбора команд
func selectableCommands(rootCmd *cobra.Command) []string {
	var names []string
	for _, c := range rootCmd.Commands() {
		switch c.Name() {
		case "help", "completion", "menu":
			continue
		}
		names = append(names, c.Name())
	}
	return names
}

// buildArgs собирает список аргументов для запуска выбранной команды
func buildArgs(rootCmd *cobra.Command, name string) ([]string, error) {
	cmd, _, err := rootCmd.Find([]string{name})
	if err != nil {
		return nil, err
	}

	args, err := collectArgs(cmd)
	if err != nil {
		return nil, err
	}

	return append([]string{name}, args...), nil
}

// collectArgs собирает аргументы для команды (подкоманды, валидные значения, произвольные аргументы)
func collectArgs(cmd *cobra.Command) ([]string, error) {
	// если у команды есть подкоманды (например migrate) — выбираем одну из них
	subs := selectableCommands(cmd)
	if len(subs) > 0 {
		selected, err := pterm.DefaultInteractiveSelect.
			WithOptions(subs).
			WithDefaultText(fmt.Sprintf("Выберите подкоманду %s", cmd.Name())).
			WithMaxHeight(8).
			Show()
		if err != nil {
			return nil, err
		}

		childArgs, err := collectArgs(childCommand(cmd, selected))
		if err != nil {
			return nil, err
		}

		return append([]string{selected}, childArgs...), nil
	}

	// команда имеет фиксированный набор валидных аргументов (например months)
	if len(cmd.ValidArgs) > 0 {
		selected, err := pterm.DefaultInteractiveSelect.
			WithOptions(cmd.ValidArgs).
			WithDefaultText(fmt.Sprintf("Выберите значение для %s", cmd.Name())).
			WithMaxHeight(8).
			Show()
		if err != nil {
			return nil, err
		}

		return []string{selected}, nil
	}

	// команда требует аргументы — запрашиваем их вручную
	if cmd.Args != nil && cmd.Args(cmd, []string{}) != nil {
		input, err := pterm.DefaultInteractiveTextInput.
			WithDefaultText(fmt.Sprintf("Введите аргументы для %s (через пробел)", cmd.Name())).
			Show()
		if err != nil {
			return nil, err
		}

		return strings.Fields(input), nil
	}

	return nil, nil
}

// childCommand возвращает подкоманду по имени
func childCommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, c := range cmd.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return cmd
}
