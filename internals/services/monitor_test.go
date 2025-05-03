package services

import (
	"testing"

	"github.com/ViniMendes2515/price-crawler/internals/models"
)

func TestAgruparSite(t *testing.T) {
	produtos := []models.ProductInfo{
		{Url: "https://kabum.com.br/produto1"},
		{Url: "https://kabum.com.br/produto2"},
		{Url: "https://amazon.com.br/produto3"},
		{Url: "https://mercadolivre.com.br/produto4"},
		{Url: "https://mercadolivre.com.br/produto5"},
	}

	grupos := agrupaSite(produtos)

	if len(grupos["kabum.com.br"]) != 2 {
		t.Errorf("Esperado 2 produtos para kabum.com.br, mas obteve %d", len(grupos["kabum.com.br"]))
	}

	if condition := len(grupos["amazon.com.br"]); condition != 1 {
		t.Errorf("Esperado 1 produto para amazon.com.br, mas obteve %d", condition)
	}

	if condition := len(grupos["mercadolivre.com.br"]); condition != 2 {
		t.Errorf("Esperado 2 produtos para mercadolivre.com.br, mas obteve %d", condition)
	}
}
