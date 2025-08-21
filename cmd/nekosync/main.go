package main

import (
	"nekosync/internal/config"
	"nekosync/internal/infrastructure/database"
	"nekosync/internal/interfaces/http"
)

func main() {
	cfg := config.Load()

	db := database.Init(cfg)
	defer db.Close()

	server := http.NewServer(cfg, db)
	server.Start(":" + cfg.Port)
}
