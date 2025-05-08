package historico

import (
	"errors"
	"math"
	"vigil/database"
)

// RegistrarPreco registra o preco de um produto na tabela historico_precos
func RegistrarPreco(url string, preco float64) error {
	db := database.DB

	var urlId int
	err := db.QueryRow("SELECT id FROM urls WHERE url = ?", url).Scan(&urlId)
	if err != nil {
		return errors.New("URL nao encontrada")
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM historico_precos WHERE url_id = ? AND preco = ?", urlId, preco).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Preco ja registrado
	}

	_, err = db.Exec("INSERT INTO historico_precos (url_id, preco) VALUES (?, ?)", urlId, preco)
	return err
}

// MenorPreco retorna o menor preco registrado para uma URL
func MenorPreco(url string) (float64, error) {
	db := database.DB

	var menor float64
	err := db.QueryRow(`
		SELECT MIN(preco)
		FROM historico_precos
		WHERE url_id = (SELECT id FROM urls WHERE url = ?)
	`, url).Scan(&menor)

	if err != nil || menor == 0 {
		return 0, errors.New("nenhum preco encontrado")
	}

	return menor, nil
}

// Media retorna a media dos precos registrados para uma URL
func Media(url string) (float64, error) {
	db := database.DB

	var media float64
	err := db.QueryRow(`SELECT AVG(preco)
		FROM historico_precos
		WHERE url_id = (SELECT id FROM urls WHERE url = ?)
	`, url).Scan(&media)

	if err != nil || media == 0 {
		return 0, errors.New("nenhum preco encontrado")
	}

	return media, nil
}

// DesvioPadrao retorna o desvio padrao dos precos registrados para uma URL
func DesvioPadrao(url string) (float64, error) {
	db := database.DB

	rows, err := db.Query(`
		SELECT preco
		FROM historico_precos
		WHERE url_id = (SELECT id FROM urls WHERE url = ?)
	`, url)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var precos []float64
	var soma float64

	for rows.Next() {
		var preco float64
		if err := rows.Scan(&preco); err != nil {
			return 0, err
		}
		precos = append(precos, preco)
		soma += preco
	}

	n := float64(len(precos))
	if n == 0 {
		return 0, errors.New("nenhum preco encontrado para o calculo")
	}

	media := soma / n
	var somaQuadrados float64
	for _, preco := range precos {
		somaQuadrados += (preco - media) * (preco - media)
	}

	variancia := somaQuadrados / (n - 1)
	return math.Sqrt(variancia), nil

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
