package api

import (
	"encoding/json"
	"net/http"

	"github.com/livrasand/ethicalmetrics/internal/db"
)

type Event struct {
	EventType string `json:"evento"`
	Module    string `json:"modulo"`
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
	var e Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil || e.EventType == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(
		"INSERT INTO events (event_type, module, duration_ms) VALUES (?, ?, ?)",
		e.EventType, e.Module, e.Duration)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	rows1, _ := db.DB.Query(`SELECT module, COUNT(*) FROM events GROUP BY module`)
	var porModulo []moduloStat
	for rows1.Next() {
		var m moduloStat
		rows1.Scan(&m.Modulo, &m.Total)
		porModulo = append(porModulo, m)
	}

	rows2, _ := db.DB.Query(`SELECT strftime('%Y-%m-%d', timestamp), COUNT(*) FROM events GROUP BY 1`)
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
