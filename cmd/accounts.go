package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"money-stat/internal/services/zenmoney"
	"money-stat/internal/usecase"
	"net/http"
)

func RunAccountList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Показать счета с балансом и валютой",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		log.Println("Show accounts called")
		api := zenmoney.NewApi(&http.Client{})

		accounts := usecase.NewAccounts(api)

		accounts.GetAccounts()

		return nil
	}

	return cmd
}
