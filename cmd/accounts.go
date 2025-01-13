package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"log"
	"money-stat/internal/services/zenmoney"
	"net/http"
	"strconv"
)

func RunAccountList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Показать счета с балансом и валютой",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		log.Println("Show accounts called")
		api := zenmoney.NewApi(&http.Client{})

		result, err := api.Diff()

		if err != nil {
			log.Println(err)
		}

		tableData := pterm.TableData{
			{"Счет", "Баланс", "Валюта"},
			{" ", " ", " "},
		}

		instruments := make(map[int]string)

		for _, instrument := range result.Instrument {
			instruments[instrument.Id] = instrument.Title
		}

		var summRuble float64
		var summDollar float64

		for _, account := range result.Account {
			tableData = append(tableData, []string{account.Title, strconv.FormatFloat(account.Balance, 'f', 2, 64), instruments[account.Instrument]})
			if account.IsRuble() {
				summRuble = summRuble + account.Balance
			}
			if account.IsDollar() {
				summDollar = summDollar + account.Balance
			}
		}

		pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()

		summData := pterm.TableData{
			{
				"Итого в рублях",
				"Итого в долларах",
			},
			{" ", " "},
			{
				strconv.FormatFloat(summRuble, 'f', 2, 64),
				strconv.FormatFloat(summDollar, 'f', 2, 64),
			},
		}

		pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(summData).Render()

		return nil
	}

	return cmd
}
