package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"money-stat/cmd"
	app2 "money-stat/internal/app"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	db := app2.NewDB()
	c := app2.NewContainer(db)
	app := app2.NewApp(c)

	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(
		cmd.RunAccountList(app),
		cmd.RunMonths(app),
		cmd.RunYear(app),
		cmd.RunSync(app),
		cmd.RunCapital(app),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute root cmd: %v", err)
	}
}
