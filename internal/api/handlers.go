package api

import (
	"encoding/json"
	"net/http"
	"io"
	"github.com/livrasand/ethicalmetrics/internal/db"
	"math/rand"
	"time"
	"github.com/google/uuid"
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
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 24)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type Event struct {
	EventType string `json:"evento"`
	Module    string `json:"modulo"`
	SiteID    string `json:"site_id"`
	Duration  int    `json:"duracion_ms"`
}

type moduloStat struct {
	Modulo string `json:"modulo"`
	Total  int    `json:"total"`
}

type diaStat struct {
	Dia   string `json:"dia"`
	Total int    `json:"total"`
}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo body", http.StatusBadRequest)
		return
	}

	var e Event
	err = json.Unmarshal(body, &e)
	if err != nil {
		// Mostrar qué llegó realmente
		http.Error(w, "JSON inválido: "+err.Error()+" || Body: "+string(body), http.StatusBadRequest)
		return
	}

	if e.EventType == "" || e.SiteID == "" {
		http.Error(w, "Campos obligatorios faltantes", http.StatusBadRequest)
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

	_, err = db.DB.Exec(
		"INSERT INTO events (event_type, module, site_id, duration_ms) VALUES (?, ?, ?, ?)",
		e.EventType, e.Module, e.SiteID, e.Duration)
	if err != nil {
		http.Error(w, "Error al guardar evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	siteID := r.URL.Query().Get("site")
	token := r.URL.Query().Get("token")

	if siteID == "" || token == "" {
		http.Error(w, "Faltan parámetros", http.StatusBadRequest)
		return
	}

	// Verificar token válido para ese site
	row := db.DB.QueryRow("SELECT COUNT(*) FROM sites WHERE id = ? AND admin_token = ?", siteID, token)
	var count int
	row.Scan(&count)
	if count == 0 {
		http.Error(w, "Token inválido", http.StatusForbidden)
		return
	}

	rows1, _ := db.DB.Query(`
		SELECT module, COUNT(*) 
		FROM events 
		WHERE site_id = ? 
		GROUP BY module`, siteID)
	var porModulo []moduloStat
	for rows1.Next() {
		var m moduloStat
		rows1.Scan(&m.Modulo, &m.Total)
		porModulo = append(porModulo, m)
	}

	rows2, _ := db.DB.Query(`
		SELECT strftime('%Y-%m-%d', timestamp), COUNT(*) 
		FROM events 
		WHERE site_id = ? 
		GROUP BY 1`, siteID)
	var porDia []diaStat
	for rows2.Next() {
		var d diaStat
		rows2.Scan(&d.Dia, &d.Total)
		porDia = append(porDia, d)
	}

	resp := map[string]interface{}{
		"por_modulo": porModulo,
		"por_dia":    porDia,
	}
	json.NewEncoder(w).Encode(resp)
}
