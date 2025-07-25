package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"money-stat/internal/app"
	"money-stat/internal/usecase"
	"strconv"
	"strings"
	"time"
)

func RunCapital() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capital",
		Short: "Показать капитал за указанный год",
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

		tableData := pterm.TableData{
			{"Месяц", "Капитал"},
			{" ", " "},
		}

		capital := usecase.NewCapital(
			app.GetGlobalApp().GetContainer().GetTransactionRepository(),
			app.GetGlobalApp().GetContainer().GetAccountRepository(),
		)

		valuesSlice, err := capital.GetCapital()

		if err != nil {
			return err
		}

		for _, row := range valuesSlice {
			if !strings.HasPrefix(row.Month, strconv.Itoa(selectYear)+"-") {
				continue
			}
			tableData = append(
				tableData,
				[]string{
					row.Month,
					strconv.FormatFloat(row.Balance, 'f', 2, 64),
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
