package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"money-stat/internal/app"
	"money-stat/internal/usecase"
	"strconv"
	"time"
)

func RunYear() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "year",
		Short: "Показать таблицу доходов и расходов  за указанный год",
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

		year := usecase.NewYear(app.GetGlobalApp().GetContainer().GetTransactionRepository())

		valuesSlice := year.GetYearStat(selectYear)

		tableData := pterm.TableData{
			{"Месяц", "Доход", "Расход", "Чистыми"},
			{" ", " ", " ", " "},
		}

		var prevDiff string
		var prevInDiff string
		for index, row := range valuesSlice {
			prevIndex := 0
			if index > 0 {
				prevIndex = index - 1
			}
			diff := row.Income - row.OutCome
			var diffStr string
			if diff < 0 {
				diffStr = pterm.FgRed.Sprint(strconv.FormatFloat(diff, 'f', 2, 64))
			} else {
				diffStr = pterm.FgGreen.Sprint(strconv.FormatFloat(diff, 'f', 2, 64))
			}
			if row.OutCome > valuesSlice[prevIndex].OutCome {
				prevDiff = "↑"
			} else {
				prevDiff = "↓"
			}

			if row.Income > valuesSlice[prevIndex].Income {
				prevInDiff = "↑"
			} else {
				prevInDiff = "↓"
			}
			tableData = append(
				tableData,
				[]string{
					row.Month,
					strconv.FormatFloat(row.Income, 'f', 2, 64) + " " + prevInDiff,
					strconv.FormatFloat(row.OutCome, 'f', 2, 64) + " " + prevDiff,
					diffStr,
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
