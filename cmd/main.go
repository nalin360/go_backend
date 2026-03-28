package main

import (
	// "log"
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"wealtharena.in/api/internal/env"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("DATABASE_URL", ""),
		},
	}


	// dotenv
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// slog
	slog.SetDefault(logger)

	// database
	connPool, err := pgxpool.New(ctx, dbURL)

	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database")
	defer connPool.Close()
	

	api := application{
		config: cfg,
		db:connPool,
	}
	if err := api.run(api.mount()); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
