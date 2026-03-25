package main

import (
	// "log"
	"log/slog"
	"os"
)


func main() {

	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,

	}

	logger := slog.New(slog.NewTextHandler(os.Stdout , nil ))

	// slog  	
	slog.SetDefault(logger)
	if err:= api.run(api.mount()); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
