package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"log"
	"money-stat/internal/app"
	"money-stat/internal/usecase"
	"strconv"
)

func RunAccountList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Показать счета с балансом и валютой",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		log.Println("Show accounts called")

		accounts := usecase.NewAccounts(app.GetGlobalApp().GetContainer().GetAccountRepository())

		stat := accounts.GetAccounts()

		tableData := pterm.TableData{
			{"Счет", "Баланс", "Валюта"},
			{" ", " ", " "},
		}

		for _, account := range stat.Accounts {
			tableData = append(tableData, []string{account.Account, account.Balance, account.Currency})
		}

		errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
		if errTable != nil {
			fmt.Println(errTable)
		}

		summData := pterm.TableData{
			{
				"Итого в рублях",
				"Итого в долларах",
				"Общая сумма в рублях",
			},
			{" ", " "},
			{
				strconv.FormatFloat(stat.SummRuble, 'f', 2, 64),
				strconv.FormatFloat(stat.SummDollar, 'f', 2, 64),
				strconv.FormatFloat(stat.SummRuble+(stat.SummDollar*stat.RateDollar), 'f', 2, 64),
			},
		}

		errSummTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(summData).Render()
		if errSummTable != nil {
			fmt.Println(errSummTable)
		}

		return nil
	}

	return cmd
}
