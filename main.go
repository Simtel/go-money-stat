package money_stat

import (
	"github.com/spf13/cobra"
	"money-stat/cmd"
)

func main() {

	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(
		cmd.RunStat(),
	)
}
