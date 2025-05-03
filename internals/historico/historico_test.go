package historico

import "testing"

func TestRegistrarPreco(t *testing.T) {
	url := "https://example.com/produto1"
	preco := 100.0

	RegistrarPreco(url, preco)

	if len(dados[url]) != 1 {
		t.Errorf("Esperado 1 preco registrado, mas obteve %d", len(dados[url]))
	}
}

func TestMenorPreco(t *testing.T) {
	url := "kabum.com.br/produto1"
	dados[url] = []float64{100.0, 90.0, 80.0}

	menor, err := MenorPreco(url)
	if err != nil {
		t.Errorf("Erro ao obter menor preco: %v", err)
	}

	if menor != 80.0 {
		t.Errorf("Esperado menor preco 80.0, mas obteve %.2f", menor)
	}
}

func TestMedia(t *testing.T) {
	url := "kabum.com.br/produto1"
	dados[url] = []float64{100.0, 90.0, 80.0}

	media, err := Media(url)

	if err != nil {
		t.Errorf("Erro ao obter media: %v", err)
	}

	if media != 90.0 {
		t.Errorf("Esperado media 90.0, mas obteve %.2f", media)
	}
}

func TestDesvioPadrao(t *testing.T) {
	url := "kabum.com.br/produto1"
	dados[url] = []float64{100.0, 90.0, 80.0}

	desvio, err := DesvioPadrao(url)

	if err != nil {
		t.Errorf("Erro ao obter desvio padrao: %v", err)
	}

	if desvio != 10 {
		t.Errorf("Esperado desvio padrao 8.16, mas obteve %.2f", desvio)
	}

}

func TestDetectarNovaMinima(t *testing.T) {
	url := "kabum.com.br/produto1"
	dados[url] = []float64{100.0, 90.0, 80.0}

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
	url := "kabum.com.br/produto1"
	dados[url] = []float64{120.0, 140.0, 80.0}

	resultado, err := DetectarPromocao(url, 85.0, 0.6)
	if err != nil {
		t.Fatalf("Erro: %v", err)
	}

	if !resultado {
		t.Errorf("Esperado promocao detectada, mas obteve %v", resultado)
	}

	resultado, _ = DetectarPromocao(url, 150, 0.6)
	if resultado {
		t.Errorf("Erro ao detectar promocao: %v", err)
	}
}
