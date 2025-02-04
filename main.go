package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"money-stat/cmd"
	app2 "money-stat/internal/app"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	db := app2.NewDB()
	c := app2.NewContainer(db)
	app := app2.NewApp(c)

	app2.SetGlobalApp(app)
}

func main() {

	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(
		cmd.RunAccountList(),
		cmd.RunMonths(),
		cmd.RunYear(),
		cmd.RunSync(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute root cmd: %v", err)

		return
	}
}
