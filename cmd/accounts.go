package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"money-stat/internal/services/zenmoney"
	"net/http"
)

func RunAccountList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Show list accounts",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		log.Println("Show accounts called")
		api := zenmoney.NewApi(&http.Client{})

		resp, err := api.Diff()

		if err != nil {
			log.Println(err)
		}

		log.Println(resp)

		return nil
	}

	return cmd
}
