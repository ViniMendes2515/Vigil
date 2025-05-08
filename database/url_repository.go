package database

import (
	"errors"
	"fmt"
	"strings"
)

// AddUrl adiciona uma nova URL ao banco de dados
func AddUrl(site string, url string, name string, precoInicial float64, precoLimite float64) error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("erro ao iniciar transação")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var validate int
	err = DB.QueryRow(`SELECT COUNT(*) FROM urls WHERE url = ?`, url).Scan(&validate)
	if err != nil {
		return err
	}

	if validate > 0 {
		return fmt.Errorf("URL ja existe no banco")
	}

	res, err := DB.Exec(`INSERT INTO urls (site, url, name, preco_limite) VALUES (?, ?, ?, ?)`, strings.ToLower(site), url, name, precoLimite)
	if err != nil {
		return errors.New("erro ao adicionar URL")
	}

	if precoInicial > 0 {
		lastId, err := res.LastInsertId()
		if err != nil {
			return errors.New("erro ao obter ID do produto")
		}
		_, err = DB.Exec(`INSERT INTO historico_precos (url_id, preco) VALUES (?, ?)`, lastId, precoInicial)
		if err != nil {
			return errors.New("erro ao adicionar preco inicial")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.New("erro ao confirmar transação")
	}

	return nil
}

// GetUrls retorna todas as URLs do banco de dados
func GetUrls() (map[string]int, error) {
	rows, err := DB.Query(`SELECT id, url FROM urls ORDER BY id`)
	if err != nil {
		return nil, err
	}

	urls := map[string]int{}

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
func GetSiteUrls(host string) ([]string, error) {
	rows, err := DB.Query(`SELECT url FROM urls WHERE site = ?`, host)
	if err != nil {
		return nil, err
	}

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if len(urls) == 0 {
		fmt.Printf("⚠️- Nenhuma URL de %s foi encontrada para monitoramento\n", host)
	}

	return urls, nil
}

// RemoveUrl remove a URL do banco de dados
func RemoveUrl(url string) error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("erro ao iniciar transação")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = DB.Exec(`DELETE FROM historico_precos WHERE url_id = (SELECT url_id FROM urls WHERE url = ? )`, url)
	if err != nil {
		return errors.New("erro ao remover historico de precos")
	}

	_, err = DB.Exec(`DELETE FROM urls WHERE url = ?`, url)
	if err != nil {
		return errors.New("erro ao remover URL")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("erro ao confirmar transação")
	}

	return nil
}

// RemoveAllUrls remove todas as URLs do banco de dados
func RemoveAllUrls() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("erro ao iniciar transação")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = DB.Exec(`DELETE FROM historico_precos`)
	if err != nil {
		return errors.New("erro ao remover historico de precos")
	}

	_, err = DB.Exec(`DELETE FROM urls`)
	if err != nil {
		return errors.New("erro ao remover URLs")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("erro ao confirmar transação")
	}

	return nil
}

// RemoveUrlById remove a URL do banco de dados com base no ID
func RemoveUrlById(id int) error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("erro ao iniciar transação")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = DB.Exec(`DELETE FROM historico_precos WHERE url_id = ?`, id)
	if err != nil {
		return errors.New("erro ao remover historico de precos")
	}

	_, err = DB.Exec(`DELETE FROM urls WHERE id = ?`, id)
	if err != nil {
		return errors.New("erro ao remover URL")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("erro ao confirmar transação")
	}

	return nil
}
