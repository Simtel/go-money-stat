package dynamics

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"money-stat/internal/app"
	"money-stat/internal/usecase"
	"strconv"
	"time"
)

func Run(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dynamics [год]",
		Short: "Показать динамику доходов и расходов по месяцам за указанный год",
		Long: `Показывает помесячную таблицу доходов и расходов с абсолютными и 
процентными изменениями относительно предыдущего месяца.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err
			}

			selectYear, errYear := strconv.Atoi(args[0])

			if errYear != nil {
				return errYear
			}

			if selectYear > 2020 && selectYear <= time.Now().Year() {
				return nil
			}
			return fmt.Errorf("указан неверный год: %s", args[0])
		},
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		selectYear, _ := strconv.Atoi(args[0])

		dynamics := usecase.NewDynamics(app.GetContainer().GetTransactionRepository())

		valuesSlice, err := dynamics.GetDynamics(selectYear)
		if err != nil {
			return err
		}

		if len(valuesSlice) == 0 {
			pterm.Info.Println("Нет данных за указанный год")
			return nil
		}

		tableData := pterm.TableData{
			{"Месяц", "Доход", "Δ Доход", "Δ%", "Расход", "Δ Расход", "Δ%", "Чистыми"},
			{" ", " ", " ", " ", " ", " ", " ", " "},
		}

		for _, row := range valuesSlice {
			net := row.Income - row.Outcome
			var netStr string
			if net < 0 {
				netStr = pterm.FgRed.Sprint(strconv.FormatFloat(net, 'f', 2, 64))
			} else {
				netStr = pterm.FgGreen.Sprint(strconv.FormatFloat(net, 'f', 2, 64))
			}

			incChangeStr := formatChange(row.IncomeChange)
			outChangeStr := formatChange(row.OutcomeChange)
			incPctStr := formatPercentChange(row.IncomeChangePct)
			outPctStr := formatPercentChange(row.OutcomeChangePct)

			tableData = append(
				tableData,
				[]string{
					row.Month,
					strconv.FormatFloat(row.Income, 'f', 2, 64),
					incChangeStr,
					incPctStr,
					strconv.FormatFloat(row.Outcome, 'f', 2, 64),
					outChangeStr,
					outPctStr,
					netStr,
				},
			)
		}

		errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
		if errTable != nil {
			fmt.Println(errTable)
		}

		return nil
	}

	return cmd
}

// formatChange форматирует абсолютное изменение с цветом
func formatChange(change float64) string {
	if change == 0 {
		return "  0.00"
	}
	if change > 0 {
		return pterm.FgRed.Sprint("+" + strconv.FormatFloat(change, 'f', 2, 64))
	}
	return pterm.FgGreen.Sprint(strconv.FormatFloat(change, 'f', 2, 64))
}

// formatPercentChange форматирует процентное изменение с цветом
func formatPercentChange(pct float64) string {
	if pct == 0 {
		return "   0%"
	}
	if pct > 0 {
		return pterm.FgRed.Sprint("+" + strconv.FormatFloat(pct, 'f', 0, 64) + "%")
	}
	return pterm.FgGreen.Sprint(strconv.FormatFloat(pct, 'f', 0, 64) + "%")
}
