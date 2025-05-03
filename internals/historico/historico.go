package historico

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"sync"
)

var mu sync.Mutex

const historicoPath = "../historico_precos.json"

var dados = make(map[string][]float64)

// Carregar o historico de precos do arquivo JSON
func Carregar() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(historicoPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(&dados)
}

// Salvar o historico de precos no arquivo JSON
func Salvar() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create(historicoPath)
	if err != nil {
		return err
	}

	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	return enc.Encode(dados)
}

// RegistrarPreco registra o preco de um produto na memoria do programa
func RegistrarPreco(url string, preco float64) {
	mu.Lock()
	defer mu.Unlock()

	for _, p := range dados[url] {
		if p == preco {
			return // Preco ja registrado
		}
	}

	dados[url] = append(dados[url], preco)
}

// MenorPreco retorna o menor preco registrado para uma URL
func MenorPreco(url string) (float64, error) {
	mu.Lock()
	defer mu.Unlock()

	precos, existe := dados[url]
	if !existe || len(precos) == 0 {
		return 0, errors.New("nenhum preco encontrado")
	}

	menor := precos[0]
	for _, preco := range precos {
		if preco < menor {
			menor = preco
		}
	}

	return menor, nil
}

// Media retorna a media dos precos registrados para uma URL
func Media(url string) (float64, error) {
	mu.Lock()
	defer mu.Unlock()

	precos, existe := dados[url]
	if !existe || len(precos) == 0 {
		return 0, errors.New("nenhum preco encontrado")
	}

	soma := 0.0
	for _, preco := range precos {
		soma += preco
	}

	media := soma / float64(len(precos))

	return media, nil
}

// DesvioPadrao retorna o desvio padrao dos precos registrados para uma URL
func DesvioPadrao(url string) (float64, error) {
	mu.Lock()
	defer mu.Unlock()

	precos, existe := dados[url]
	if !existe || len(precos) == 0 {
		return 0, errors.New("nenhum preco encontrado")
	}

	var soma, media float64
	for _, preco := range precos {
		soma += preco
	}
	media = soma / float64(len(precos))

	var somaQuadrados float64
	for _, preco := range precos {
		somaQuadrados += (preco - media) * (preco - media)
	}

	desvio := somaQuadrados / float64(len(precos)-1)
	return math.Sqrt(desvio), nil
}

// DetectarNovaMinima verifica se o preco atual e menor que o menor preco registrado
func DetectarNovaMinima(url string, precoAtual float64) (bool, error) {
	menor, err := MenorPreco(url)
	if err != nil {
		return true, err
	}

	return precoAtual < menor, nil
}

// DetectarPromocao verifica se o preco atual e menor que a media menos o desvio padrao
func DetectarPromocao(url string, precoAtual float64, fatorDesvio float64) (bool, error) {
	minimo, _ := DetectarNovaMinima(url, precoAtual)
	if minimo {
		return true, nil
	}

	media, err := Media(url)
	if err != nil {
		return true, err // Caso nao tenha preco registrado para URL
	}

	desvio, err := DesvioPadrao(url)
	if err != nil {
		return true, err // Caso nao tenha preco registrado para URL
	}

	limite := media - (desvio * fatorDesvio)
	return precoAtual < limite, nil

}
