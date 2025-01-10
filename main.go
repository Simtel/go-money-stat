package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"money-stat/cmd"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(
		cmd.RunAccountList(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute root cmd: %v", err)

		return
	}
}
