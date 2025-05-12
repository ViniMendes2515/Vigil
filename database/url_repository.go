package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type UrlRepository interface {
	AddUrl(ctx context.Context, site, url, name string, precoInicial, precoLimite float64) error
	GetUrls(ctx context.Context) (map[string]int, error)
	GetSiteUrls(ctx context.Context, host string) ([]string, error)
	RemoveUrl(ctx context.Context, url string) error
	RemoveAllUrls(ctx context.Context) error
	RemoveUrlById(ctx context.Context, id uint) error
}

type PostgresUrlRepo struct{}

// AddUrl adiciona uma nova URL ao banco de dados
func (r *PostgresUrlRepo) AddUrl(ctx context.Context, site string, url string, name string, precoInicial float64, precoLimite float64) error {
	tx, err := DB.Begin(ctx)
	if err != nil {
		return errors.New("erro ao iniciar transação")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var validate int
	err = DB.QueryRow(ctx, `SELECT COUNT(*) FROM urls WHERE url = $1`, url).Scan(&validate)
	if err != nil {
		return err
	}

	if validate > 0 {
		return fmt.Errorf("URL já existe no banco")
	}

	var lastId int64
	err = tx.QueryRow(ctx,
		`INSERT INTO urls (site, url, name, preco_limite) VALUES ($1, $2, $3, $4) RETURNING id`,
		strings.ToLower(site), url, name, precoLimite).Scan(&lastId)

	if err != nil {
		return errors.New("erro ao adicionar URL")
	}

	if precoInicial > 0 {
		_, err = tx.Exec(ctx,
			`INSERT INTO historico_precos (url_id, preco) VALUES ($1, $2)`,
			lastId, precoInicial)

		if err != nil {
			return errors.New("erro ao adicionar preco inicial")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New("erro ao confirmar transação")
	}

	return nil
}

// GetUrls retorna todas as URLs do banco de dados
func (r *PostgresUrlRepo) GetUrls(ctx context.Context) (map[string]int, error) {
	rows, err := DB.Query(ctx, `SELECT id, url FROM urls ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make(map[string]int)
	for rows.Next() {
		var url string
		var id uint
		if err := rows.Scan(&id, &url); err != nil {
			return nil, err
		}
		urls[url] = int(id)
	}

	return urls, nil
}

// GetSiteUrls retorna todas as URLs de um site específico
func (r *PostgresUrlRepo) GetSiteUrls(ctx context.Context, host string) ([]string, error) {
	rows, err := DB.Query(ctx, `SELECT url FROM urls WHERE site = $1`, host)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if len(urls) == 0 {
		fmt.Printf("⚠️ Nenhuma URL de %s foi encontrada para monitoramento\n", host)
	}

	return urls, nil
}

// RemoveUrl remove uma URL do banco de dados
func (r *PostgresUrlRepo) RemoveUrl(ctx context.Context, url string) error {
	tx, err := DB.Begin(ctx)
	if err != nil {
		return errors.New("erro ao iniciar transação")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx,
		`DELETE FROM historico_precos WHERE url_id = (SELECT id FROM urls WHERE url = $1)`, url)
	if err != nil {
		return errors.New("erro ao remover historico de precos")
	}

	_, err = tx.Exec(ctx, `DELETE FROM urls WHERE url = $1`, url)
	if err != nil {
		return errors.New("erro ao remover URL")
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New("erro ao confirmar transação")
	}

	return nil
}

// RemoveAllUrls remove todas as URLs do banco de dados
func (r *PostgresUrlRepo) RemoveAllUrls(ctx context.Context) error {
	tx, err := DB.Begin(ctx)
	if err != nil {
		return errors.New("erro ao iniciar transação")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `DELETE FROM historico_precos`)
	if err != nil {
		return errors.New("erro ao remover historico de precos")
	}

	_, err = tx.Exec(ctx, `DELETE FROM urls`)
	if err != nil {
		return errors.New("erro ao remover URLs")
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New("erro ao confirmar transação")
	}

	return nil
}

// RemoveUrlById remove uma URL do banco de dados com base no ID
func (r *PostgresUrlRepo) RemoveUrlById(ctx context.Context, id uint) error {
	tx, err := DB.Begin(ctx)
	if err != nil {
		return errors.New("erro ao iniciar transação")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `DELETE FROM historico_precos WHERE url_id = $1`, id)
	if err != nil {
		return errors.New("erro ao remover historico de precos")
	}

	_, err = tx.Exec(ctx, `DELETE FROM urls WHERE id = $1`, id)
	if err != nil {
		return errors.New("erro ao remover URL")
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New("erro ao confirmar transação")
	}

	return nil
}
