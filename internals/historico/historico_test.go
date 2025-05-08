package historico

import (
	"testing"
	"vigil/database"
	"vigil/tests"
)

// Helper para criar uma URL no banco
func criarURL(t *testing.T, url string) {
	_, err := database.DB.Exec(`INSERT INTO urls (site, url, preco_limite) VALUES (?, ?, ?)`, "kabum", url, 0.0)
	if err != nil {
		t.Fatalf("Erro ao inserir URL de teste: %v", err)
	}
}

func TestRegistrarPreco(t *testing.T) {
	tests.SetupTestDB()
	url := "https://example.com/produto1"
	criarURL(t, url)

	err := RegistrarPreco(url, 100.0)
	if err != nil {
		t.Fatalf("Erro ao registrar preco: %v", err)
	}

	var count int
	err = database.DB.QueryRow(`SELECT COUNT(*) FROM historico_precos`).Scan(&count)
	if err != nil {
		t.Fatalf("Erro ao contar precos: %v", err)
	}
	if count != 1 {
		t.Errorf("Esperado 1 preco registrado, obteve %d", count)
	}
}

func TestMenorPreco(t *testing.T) {
	tests.SetupTestDB()
	url := "kabum.com.br/produto1"
	criarURL(t, url)

	RegistrarPreco(url, 100.0)
	RegistrarPreco(url, 90.0)
	RegistrarPreco(url, 80.0)

	menor, err := MenorPreco(url)
	if err != nil {
		t.Errorf("Erro ao obter menor preco: %v", err)
	}

	if menor != 80.0 {
		t.Errorf("Esperado menor preco 80.0, mas obteve %.2f", menor)
	}
}

func TestMedia(t *testing.T) {
	tests.SetupTestDB()
	url := "kabum.com.br/produto1"
	criarURL(t, url)

	RegistrarPreco(url, 100.0)
	RegistrarPreco(url, 90.0)
	RegistrarPreco(url, 80.0)

	media, err := Media(url)
	if err != nil {
		t.Errorf("Erro ao obter media: %v", err)
	}

	if media != 90.0 {
		t.Errorf("Esperado media 90.0, mas obteve %.2f", media)
	}
}

func TestDesvioPadrao(t *testing.T) {
	tests.SetupTestDB()
	url := "kabum.com.br/produto1"
	criarURL(t, url)

	RegistrarPreco(url, 100.0)
	RegistrarPreco(url, 90.0)
	RegistrarPreco(url, 80.0)

	desvio, err := DesvioPadrao(url)
	if err != nil {
		t.Errorf("Erro ao obter desvio padrao: %v", err)
	}

	if desvio < 8.1 || desvio > 8.2 {
		t.Errorf("Esperado desvio padrao aproximado 8.16, mas obteve %.2f", desvio)
	}
}

func TestDetectarNovaMinima(t *testing.T) {
	tests.SetupTestDB()
	url := "kabum.com.br/produto1"
	criarURL(t, url)

	RegistrarPreco(url, 100.0)
	RegistrarPreco(url, 90.0)
	RegistrarPreco(url, 80.0)

	minimo1, _ := DetectarNovaMinima(url, 75.0)
	if !minimo1 {
		t.Errorf("Esperado nova minima detectada, mas obteve %v", minimo1)
	}

	minimo2, _ := DetectarNovaMinima(url, 80.0)
	if minimo2 {
		t.Errorf("Esperado nova minima nao detectada, mas obteve %v", minimo2)
	}
}

func TestDetectarPromocao(t *testing.T) {
	tests.SetupTestDB()
	url := "kabum.com.br/produto1"
	criarURL(t, url)

	RegistrarPreco(url, 120.0)
	RegistrarPreco(url, 140.0)
	RegistrarPreco(url, 80.0)

	resultado, err := DetectarPromocao(url, 85.0, 0.6)
	if err != nil {
		t.Fatalf("Erro: %v", err)
	}
	if !resultado {
		t.Errorf("Esperado promocao detectada, mas obteve %v", resultado)
	}

	resultado, _ = DetectarPromocao(url, 150, 0.6)
	if resultado {
		t.Errorf("Erro ao detectar promocao: preco muito alto foi considerado promocao")
	}
}
