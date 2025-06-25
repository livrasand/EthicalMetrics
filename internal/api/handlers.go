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

	err := db.RDB.HSet(db.Ctx, "site:"+siteID, map[string]interface{}{
		"name":        "Sitio sin nombre",
		"admin_token": adminToken,
		"created_at":  time.Now().Format(time.RFC3339),
	}).Err()

	if err != nil {
		http.Error(w, "Error creando sitio", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"site_id":     siteID,
		"admin_token": adminToken,
		"instruccion": `<script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js?id=` + siteID + `"></script>`,
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
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verifica si el sitio existe
	exists, err := db.RDB.Exists(db.Ctx, "site:"+e.SiteID).Result()
	if err != nil || exists == 0 {
		http.Error(w, "Site ID inválido", http.StatusForbidden)
		return
	}

	// Guarda el evento como JSON en una lista
	eventJSON, _ := json.Marshal(map[string]interface{}{
		"type":        e.EventType,
		"module":      e.Module,
		"duration_ms": e.Duration,
		"timestamp":   time.Now().Format(time.RFC3339),
	})

	err = db.RDB.RPush(db.Ctx, "events:"+e.SiteID, eventJSON).Err()
	if err != nil {
		http.Error(w, "Error guardando evento", http.StatusInternalServerError)
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

	// Verifica si el token es válido
	storedToken, err := db.RDB.HGet(db.Ctx, "site:"+siteID, "admin_token").Result()
	if err != nil || storedToken != token {
		http.Error(w, "Token inválido", http.StatusForbidden)
		return
	}

	// Obtener todos los eventos del sitio
	eventsRaw, err := db.RDB.LRange(db.Ctx, "events:"+siteID, 0, -1).Result()
	if err != nil {
		http.Error(w, "Error obteniendo eventos", http.StatusInternalServerError)
		return
	}

	modCount := map[string]int{}
	dayCount := map[string]int{}

	for _, raw := range eventsRaw {
		var evt map[string]interface{}
		json.Unmarshal([]byte(raw), &evt)

		mod := evt["module"].(string)
		ts := evt["timestamp"].(string)

		// agrupar por módulo
		modCount[mod]++

		// agrupar por día
		t, _ := time.Parse(time.RFC3339, ts)
		day := t.Format("2006-01-02")
		dayCount[day]++
	}

	var porModulo []moduloStat
	for m, total := range modCount {
		porModulo = append(porModulo, moduloStat{Modulo: m, Total: total})
	}

	var porDia []diaStat
	for d, total := range dayCount {
		porDia = append(porDia, diaStat{Dia: d, Total: total})
	}

	resp := map[string]interface{}{
		"por_modulo": porModulo,
		"por_dia":    porDia,
	}
	json.NewEncoder(w).Encode(resp)
}