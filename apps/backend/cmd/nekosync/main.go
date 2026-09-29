package main

import (
	"log"

	"nekosync/internal/config"
	"nekosync/internal/infrastructure/database"
	"nekosync/internal/interfaces/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := http.NewHTTPServer(cfg, db)
	server.Logger.Fatal(server.Start(":" + cfg.Port))
}
