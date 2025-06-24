package api

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/livrasand/ethicalmetrics/internal/db"
)

func NuevoHandler(w http.ResponseWriter, r *http.Request) {
	siteID := uuid.NewString()
	adminToken := generarToken()

	_, err := db.DB.Exec(
		"INSERT INTO sites (id, name, admin_token) VALUES (?, ?, ?)",
		siteID, "Sitio sin nombre", adminToken)

	if err != nil {
		http.Error(w, "Error creando sitio", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"site_id":     siteID,
		"admin_token": adminToken,
		"instruccion": `<script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js" defer data-site-id="` + siteID + `"></script>`,
	}
	json.NewEncoder(w).Encode(resp)
}

func generarToken() string {
	b := make([]byte, 24)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

type Event struct {
	EventType string `json:"evento"`
	Module    string `json:"modulo"`
	SiteID    string `json:"site_id"`
	Duration  int    `json:"duracion_ms"`
}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	// Limita el tamaño del cuerpo a 1KB para evitar abusos
	r.Body = http.MaxBytesReader(w, r.Body, 1024)

	var e Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil || e.EventType == "" || e.SiteID == "" || len(e.EventType) > 32 || len(e.Module) > 32 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Validación adicional: duración no negativa y razonable
	if e.Duration < 0 || e.Duration > 3600000 {
		http.Error(w, "Duración inválida", http.StatusBadRequest)
		return
	}

	// Validación de formato UUID para SiteID
	if _, err := uuid.Parse(e.SiteID); err != nil {
		http.Error(w, "Site ID inválido", http.StatusBadRequest)
		return
	}

	// Verifica que el sitio exista
	row := db.DB.QueryRow("SELECT COUNT(*) FROM sites WHERE id = ?", e.SiteID)
	var count int
	row.Scan(&count)
	if count == 0 {
		http.Error(w, "Site ID inválido", http.StatusForbidden)
		return
	}

	stmt, err := db.DB.Prepare(`INSERT INTO events (site_id, module, event, duration_ms, timestamp) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.SiteID, e.Module, e.EventType, e.Duration)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	siteID := r.URL.Query().Get("site")
	token := r.URL.Query().Get("token")

	// Evita loguear tokens en producción
	if ginMode := os.Getenv("GIN_MODE"); ginMode != "release" {
		log.Println("→ Recibida solicitud /stats")
		log.Println("Site:", siteID)
		log.Println("Token:", token)
	}

	if siteID == "" || token == "" {
		http.Error(w, "Parámetros requeridos", http.StatusForbidden)
		return
	}

	var count int
	err := db.DB.QueryRow(
		"SELECT COUNT(*) FROM sites WHERE id = ? AND admin_token = ?",
		siteID, token).Scan(&count)

	if err != nil {
		log.Println("Error al consultar la DB:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		log.Println("🔒 Token inválido o no coincide")
		http.Error(w, "Token inválido", http.StatusForbidden)
		return
	}

	log.Println("✅ Token válido. Consultando datos...")

	// Estadísticas por módulo (filtrando por site_id)
	rows1, err := db.DB.Query(`SELECT module, COUNT(*) FROM events WHERE site_id = ? GROUP BY module`, siteID)
	if err != nil {
		log.Println("Error en consulta por módulo:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	defer rows1.Close()

	var porModulo []map[string]interface{}
	for rows1.Next() {
		var modulo string
		var total int
		if err := rows1.Scan(&modulo, &total); err != nil {
			log.Println("Error al escanear módulo:", err)
			continue
		}
		porModulo = append(porModulo, map[string]interface{}{
			"modulo": modulo,
			"total":  total,
		})
	}

	// Estadísticas por día (filtrando por site_id)
	rows2, err := db.DB.Query(`SELECT strftime('%Y-%m-%d', timestamp), COUNT(*) FROM events WHERE site_id = ? GROUP BY 1`, siteID)
	if err != nil {
		log.Println("Error en consulta por día:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	defer rows2.Close()

	var porDia []map[string]interface{}
	for rows2.Next() {
		var dia string
		var total int
		if err := rows2.Scan(&dia, &total); err != nil {
			log.Println("Error al escanear día:", err)
			continue
		}
		porDia = append(porDia, map[string]interface{}{
			"dia":   dia,
			"total": total,
		})
	}

	resp := map[string]interface{}{
		"por_modulo": porModulo,
		"por_dia":    porDia,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Error al codificar JSON:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
	}
}
