package cmd

import "github.com/spf13/cobra"

func RunStat() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stat",
		Short: "Show base stat",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		return nil
	}

	return cmd
}
