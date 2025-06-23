package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
	key := os.Getenv("DB_KEY")
	if key == "" {
		return fmt.Errorf("DB_KEY no definida")
	}

	var err error
	DB, err = sql.Open("sqlite3", fmt.Sprintf("file:metrics.db?_pragma_key=%s&_pragma_cipher_page_size=4096", key))
	if err != nil {
		return err
	}

	return createTable()
}

func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_type TEXT,
		module TEXT,
		duration_ms INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(query)
	return err
}
