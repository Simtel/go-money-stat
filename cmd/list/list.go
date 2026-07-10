package list

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"money-stat/internal/app"
)

func Run(_ *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Показать список всех доступных команд",
		Run: func(cmd *cobra.Command, args []string) {
			// Заголовок
			pterm.DefaultSection.Println("Доступные команды money-stat")

			// Таблица команд
			tableData := pterm.TableData{
				{"Команда", "Аргументы", "Описание"},
				{"", "", ""},
				{"sync", "--incremental / --full", "Синхронизировать данные с ZenMoney"},
				{"accounts", "", "Показать счета с балансом и валютой"},
				{"months", "current | previous", "Показать транзакции за месяц"},
				{"year", "год (например 2025)", "Отчёт доходов и расходов за год"},
				{"dynamics", "год (например 2025)", "Динамика доходов/расходов по месяцам"},
				{"capital", "год (например 2025)", "Помесячный капитал за год"},
				{"migrate", "init", "Инициализация базы данных"},
			}

			pterm.DefaultTable.
				WithHasHeader().
				WithBoxed().
				WithRowSeparator("-").
				WithData(tableData).
				Render()

			// Примеры использования
			pterm.DefaultSection.Println("Примеры использования")
			pterm.Println(pterm.FgCyan.Sprint("  money-stat sync --full"))
			pterm.Println(pterm.FgCyan.Sprint("  money-stat months current"))
			pterm.Println(pterm.FgCyan.Sprint("  money-stat year 2025"))
			pterm.Println(pterm.FgCyan.Sprint("  money-stat capital 2025"))
			pterm.Println(pterm.FgCyan.Sprint("  money-stat dynamics 2025"))
			pterm.Println(pterm.FgCyan.Sprint("  money-stat accounts"))
			pterm.Println(pterm.FgCyan.Sprint("  money-stat migrate init"))
		},
	}

	return cmd
}
