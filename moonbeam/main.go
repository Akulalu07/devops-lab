package main

import (
	"moonbeam/internal/config"
	"moonbeam/internal/database"
	"moonbeam/internal/logger"
	"moonbeam/internal/router"

	"go.uber.org/zap"
)

func main() {
	log := logger.Init()
	defer log.Sync()

	cfg := config.Load()

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}

	r := router.NewRouter(log, db)

	log.Info("Starting moonbeam service", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
