package models

type Event struct {
	ID         int    `json:"id"`
	EventType  string `json:"evento"`
	Module     string `json:"modulo"`
	DurationMs int    `json:"duracion_ms"`
	Timestamp  string `json:"timestamp"`
}
