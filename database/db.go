package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	path := filepath.Join(".", "vigil.db")

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("erro ao abrir o banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error ao conectar no banco: %v", err)
	}

	DB = db
	return createTables()
}

func createTables() error {
	urlTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NULL,
		site TEXT NOT NULL,
		url TEXT NOT NULL UNIQUE,
		preco_limite REAL NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	historicoTable := `
	CREATE TABLE IF NOT EXISTS historico_precos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url_id INTEGER NOT NULL,
		preco REAL NOT NULL,
		data DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(url_id) REFERENCES urls(id)
	);`

	_, err := DB.Exec(urlTable)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de urls: %w", err)
	}

	_, err = DB.Exec(historicoTable)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de historico de precos: %w", err)
	}

	return nil
}
