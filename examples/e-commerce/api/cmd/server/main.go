package main

import (
	"github.com/yasinatesim/vela-commerce/api/internal/config"
	"github.com/yasinatesim/vela-commerce/api/internal/database"
	"github.com/yasinatesim/vela-commerce/api/internal/logger"
	"github.com/yasinatesim/vela-commerce/api/internal/server"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Error("db open failed", "err", err)
		panic(err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Error("automigrate failed", "err", err)
		panic(err)
	}

	r := server.New(cfg, db, log)
	log.Info("listening", "port", cfg.APIPort)
	if err := r.Run(":" + cfg.APIPort); err != nil {
		log.Error("server stopped", "err", err)
		panic(err)
	}
}
