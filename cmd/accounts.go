package cmd

import (
	"github.com/spf13/cobra"
	"log"
	app2 "money-stat/internal/app"
	"money-stat/internal/usecase"
)

func RunAccountList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Показать счета с балансом и валютой",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		log.Println("Show accounts called")
		app, _ := app2.GetGlobalApp()

		accounts := usecase.NewAccounts(app.GetContainer().GetAccountRepository())

		accounts.GetAccounts()

		return nil
	}

	return cmd
}
