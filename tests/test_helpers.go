package tests

import (
	"database/sql"
	"log"
	"vigil/database"

	_ "modernc.org/sqlite"
)

// SetupTestDB configura um banco de dados SQLite em memória para testes.
func SetupTestDB() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Erro ao criar banco de teste: %v", err)
	}
	database.DB = db
	createTestTables(db)
}

// createTestTables cria as tabelas necessárias para os testes.
func createTestTables(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		site TEXT NOT NULL,
		url TEXT NOT NULL UNIQUE,
		preco_limite REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE historico_precos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url_id INTEGER NOT NULL,
		preco REAL NOT NULL,
		data DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(url_id) REFERENCES urls(id)
	);`)
	if err != nil {
		log.Fatalf("Erro ao criar tabelas de teste: %v", err)
	}
}
