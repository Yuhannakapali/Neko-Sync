package main

import (
	"nekosync/internal/config"
	"nekosync/internal/infra/db"
	"nekosync/internal/interfaces/http"
)

func main() {
	cfg := config.Load()
	
	database := db.Init(cfg)
	defer database.Close()

	echo := http.NewServer(cfg, database)
	echo.Start()
}
