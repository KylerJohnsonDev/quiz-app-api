package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load .env so it's available
	_ = godotenv.Load() // Load .env from cwd; ok if file doesn't exist

	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run(api.mount()); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
