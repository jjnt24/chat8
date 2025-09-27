package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jjnt224/chat8/pkg/config"
	"github.com/jjnt224/chat8/pkg/db"
	"github.com/jjnt224/chat8/pkg/rest"
)

func main() {
	cfg := config.Load()
	database := db.MustInit(cfg)

	r := rest.NewRouter(cfg, database)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("HTTP server listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("server error:", err)
		os.Exit(1)
	}
}
