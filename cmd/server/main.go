package main

import (
	"log"
	"net/http"
	"os"

	"github.com/livrasand/ethicalmetrics/internal/api"
	"github.com/livrasand/ethicalmetrics/internal/db"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatalf("Error al iniciar la BD: %v", err)
	}

	http.HandleFunc("/nuevo", api.NuevoHandler)
	http.HandleFunc("/stats", api.StatsHandler)
	http.HandleFunc("/track", api.TrackHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("EthicalMetrics escuchando en :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
