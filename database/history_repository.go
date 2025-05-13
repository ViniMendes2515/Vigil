package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"vigil/internals/models"
)

type PostgresPriceHistoryRepo struct{}

type PriceHistoryRepository interface {
	ListPrices(ctx context.Context) ([]string, error)
	ShowByUrl(ctx context.Context, url string) (*models.ProductDetails, error)
	ShowByID(ctx context.Context, id int) (*models.ProductDetails, error)
}

// ListPrices retorna uma lista de preços dos produtos monitorados, agrupados por URL.
func (r *PostgresPriceHistoryRepo) ListPrices(ctx context.Context) ([]string, error) {
	rows, err := DB.Query(ctx, `
        SELECT u.url, hp.preco, u.id
        FROM historico_precos hp
        JOIN urls u ON hp.url_id = u.id
        ORDER BY u.url, hp.id
    `)
	if err != nil {
		return nil, errors.New("erro ao listar precos")
	}
	defer rows.Close()

	precosUrls := make(map[string][]float64)
	urlIDs := make(map[string]uint)

	for rows.Next() {
		var url string
		var preco float64
		var id uint
		if err := rows.Scan(&url, &preco, &id); err != nil {
			return nil, errors.New("erro ao escanear precos")
		}
		urlIDs[url] = id
		precosUrls[url] = append(precosUrls[url], preco)
	}

	var result []string
	for url, precos := range precosUrls {
		var precosFormat []string
		for _, p := range precos {
			precosFormat = append(precosFormat, fmt.Sprintf("R$ %.2f", p))
		}
		linha := fmt.Sprintf("- ID do produto: %d \n- URL: %s \n- Preços: %s\n", urlIDs[url], url, strings.Join(precosFormat, " | "))
		result = append(result, linha)
	}

	return result, nil
}

// getProduct busca os detalhes de um produto com base em uma cláusula WHERE e argumentos fornecidos
func getProduct(ctx context.Context, whereClause string, arg any) (*models.ProductDetails, error) {
	var detalhes models.ProductDetails

	query := `
        SELECT 
            u.id,
            u.url,  
            u.site, 
            u.name,
            u.preco_limite, 
            MIN(hp.preco) as preco_minimo, 
            MAX(hp.preco) as preco_maximo, 
            AVG(hp.preco) as preco_medio,
            (
                SELECT hp2.preco
                FROM historico_precos hp2
                WHERE hp2.url_id = u.id
                ORDER BY hp2.data DESC
                LIMIT 1
            ) AS preco_atual,
            COUNT(hp.id) as total_coletas,
            TO_CHAR(MAX(hp.data), 'DD/MM/YYYY') AS ultima_verificacao
        FROM historico_precos hp
        JOIN urls u ON hp.url_id = u.id
        WHERE ` + whereClause + `
        GROUP BY u.id`

	row := DB.QueryRow(ctx, query, arg)

	var id uint

	err := row.Scan(&id, &detalhes.Url, &detalhes.Site, &detalhes.Nome, &detalhes.PrecoLimite, &detalhes.PrecoMinimo, &detalhes.PrecoMaximo, &detalhes.PrecoMedio, &detalhes.PrecoAtual, &detalhes.TotalColetas, &detalhes.UltimaVerificacao)
	if err != nil {
		return nil, errors.New("erro ao buscar dados do produto")
	}

	rows, err := DB.Query(ctx, `
        SELECT preco, TO_CHAR(data, 'DD/MM/YYYY') AS data
        FROM historico_precos 
        WHERE url_id = $1 
        ORDER BY data DESC 
        LIMIT 5
    `, id)

	if err != nil {
		return nil, errors.New("erro ao buscar historico de precos")
	}
	defer rows.Close()

	for rows.Next() {
		var h models.HistoricoPreco
		if err := rows.Scan(&h.Preco, &h.Data); err != nil {
			return nil, errors.New("erro ao escanear precos")
		}
		detalhes.HistoricoRecentes = append(detalhes.HistoricoRecentes, h)
	}

	return &detalhes, nil
}

// ShowByUrl busca os detalhes de um produto com base na URL fornecida
func (r *PostgresPriceHistoryRepo) ShowByUrl(ctx context.Context, url string) (*models.ProductDetails, error) {
	return getProduct(ctx, "u.url = $1", url)
}

// ShowByID busca os detalhes de um produto com base no ID fornecido
func (r *PostgresPriceHistoryRepo) ShowByID(ctx context.Context, id int) (*models.ProductDetails, error) {
	return getProduct(ctx, "u.id = $1", id)
}
