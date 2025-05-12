package database

import (
	"context"
	"fmt"
	"time"
	"vigil/config"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitDB() error {
	cfg := config.Load()
	connStr := cfg.DatabaseUrl

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return fmt.Errorf("erro ao abrir o banco: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("erro ao conectar no banco: %w", err)
	}

	DB = db
	return nil
}
