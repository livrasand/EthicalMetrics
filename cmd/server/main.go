package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/livrasand/ethicalmetrics/internal/api"
	"github.com/livrasand/ethicalmetrics/internal/db"

	"github.com/joho/godotenv"
)

func downloadGeoLiteDB(path string) error {
	const url = "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

func main() {
	_ = godotenv.Load()

	err := db.Init()
	if err != nil {
		log.Fatalf("Error al iniciar la BD: %v", err)
	}

	geoIPPath := os.Getenv("GEOIP_PATH")
	if geoIPPath == "" {
		geoIPPath = "./GeoLite2-City.mmdb"
	}
	// Si no existe, descargar automáticamente
	if _, err := os.Stat(geoIPPath); os.IsNotExist(err) {
		log.Println("GeoLite2-City.mmdb no encontrado, descargando...")
		if err := downloadGeoLiteDB(geoIPPath); err != nil {
			log.Fatalf("No se pudo descargar GeoLite2-City.mmdb: %v", err)
		}
		log.Println("GeoLite2-City.mmdb descargado correctamente.")
	}
	err = api.InitGeoIP(geoIPPath)
	if err != nil {
		log.Fatalf("Error al iniciar GeoIP: %v", err)
	}

	http.HandleFunc("/nuevo", api.NuevoHandler)
	http.HandleFunc("/stats", withCORS(api.StatsHandler))
	http.Handle("/track", api.RateLimitMiddleware(http.HandlerFunc(api.TrackHandler)))
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("EthicalMetrics escuchando en :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
