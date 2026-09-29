package main

import (
	"log"

	"nekosync/internal/app"
	"nekosync/internal/platform/config"
	"nekosync/internal/platform/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := app.NewServer(cfg, db)
	server.Logger.Fatal(server.Start(":" + cfg.Port))
}
