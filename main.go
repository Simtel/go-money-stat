package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"money-stat/cmd/accounts"
	"money-stat/cmd/capital"
	"money-stat/cmd/migrate"
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

	// Проверка обязательных переменных окружения
	if os.Getenv("ZENMONEY_TOKEN") == "" {
		log.Fatal("ZENMONEY_TOKEN environment variable is required")
	}
}

func main() {
	db := app2.NewDB()
	c := app2.NewContainer(db)
	app := app2.NewApp(c)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel() // Освобождаем ресурсы контекста

	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(
		accounts.Run(app),
		months.Run(app),
		year.Run(app),
		sync.Run(app),
		capital.Run(app),
		migrate.Run(app),
	)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("Программа завершена пользователем")
			os.Exit(0)
		}
		log.Fatal(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute root cmd: %v", err)
	}
}
