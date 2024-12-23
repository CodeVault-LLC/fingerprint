package main

import (
	"fmt"
	"os"

	"github.com/codevault-llc/fingerprint/cmd/server"
	"github.com/codevault-llc/fingerprint/config"
	"github.com/codevault-llc/fingerprint/internal/database"
	"github.com/codevault-llc/fingerprint/internal/updater"
	"github.com/codevault-llc/fingerprint/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		fmt.Println("Failed to initialize logger")
		os.Exit(1)
	}

	config, err := config.NewInternalConfig()
	if err != nil {
		log.Error("Error loading config", zap.Error(err))
		os.Exit(1)
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Error("Error connecting to database %v", zap.Error(err))
		os.Exit(1)
	}

	_ = updater.NewUpdater(db)

	if err := server.StartServer(config, db); err != nil {
		log.Error("Error starting server", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Server started successfully")
}
