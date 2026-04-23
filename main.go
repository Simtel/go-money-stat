package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"money-stat/cmd/accounts"
	"money-stat/cmd/capital"
	"money-stat/cmd/months"
	"money-stat/cmd/sync"
	"money-stat/cmd/year"
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
		accounts.Run(app),
		months.Run(app),
		year.Run(app),
		sync.Run(app),
		capital.Run(app),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute root cmd: %v", err)
	}
}
