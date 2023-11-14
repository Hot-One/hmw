package main

import (
	"log"
	"net/http"

	"app/api"
	"app/config"
	"app/storage/postgres"
)

func main() {
	cfg := config.Load()

	pgconn, err := postgres.NewConnectionPostgres(&cfg)
	if err != nil {
		panic("postgres no connection: " + err.Error())
	}

	api.NewApi(&cfg, pgconn)

	log.Println("Listening...", cfg.ServerHost+cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.ServerHost+cfg.HTTPPort, nil); err != nil {
		panic("Server no run:" + err.Error())
	}
}
