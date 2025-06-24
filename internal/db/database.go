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

	// Si el archivo existe y no está cifrado, lanzar error para evitar sobrescribirlo sin cifrado
	if _, err := os.Stat("metrics.db"); err == nil {
		// Intentar abrir con clave, si falla, advertir y salir
		tmpDB, err := sql.Open("sqlite3", fmt.Sprintf("file:metrics.db?_pragma_key=%s", key))
		if err == nil {
			defer tmpDB.Close()
			_, err = tmpDB.Exec("SELECT count(*) FROM sqlite_master;")
			if err != nil {
				return fmt.Errorf("el archivo metrics.db ya existe y no está cifrado con SQLCipher, elimínalo manualmente para continuar")
			}
		}
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
<<<<<<< HEAD
	CREATE TABLE IF NOT EXISTS sites (
		id TEXT PRIMARY KEY,         -- site_id (UUID)
		name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		admin_token TEXT             -- clave secreta para el dashboard
	);

=======
>>>>>>> 97bd8ba (Agregar implementación inicial de EthicalMetrics con soporte para SQLCipher y manejo de eventos)
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_type TEXT,
		module TEXT,
<<<<<<< HEAD
		site_id TEXT,               -- se vincula con la tabla sites
		duration_ms INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(site_id) REFERENCES sites(id)
	);`

	_, err := DB.Exec(query)
	return err
}
<<<<<<< HEAD
=======
		duration_ms INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(query)
	return err
}
>>>>>>> 97bd8ba (Agregar implementación inicial de EthicalMetrics con soporte para SQLCipher y manejo de eventos)
=======
>>>>>>> d5c0bdb (Agregar archivo .gitignore para excluir metrics.db; mejorar la inicialización de la base de datos con validaciones de cifrado; actualizar el script de seguimiento para incluir site_id; añadir script en nuevo.html para la creación de cuentas.)
