package cmd

import (
	"github.com/spf13/cobra"
	"money-stat/internal/usecase"
)

func RunSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Синхронизировать данные",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		sync := usecase.Sync{}

		sync.FullSync()

		return nil
	}

	return cmd
}
