package crawler

import "vigil/internals/models"

type CrawlerFunc func([]string) ([]models.ProductInfo, error)
type NameFetcher func(string) (string, error)

var (
	registry     = make(map[string]CrawlerFunc)
	nameRegistry = make(map[string]NameFetcher)
)

// Register registra uma função de crawler para um site específico
func Register(site string, fn CrawlerFunc) {
	registry[site] = fn
}

// GetRegistered retorna todas as funções de crawler registradas
func GetRegistered() map[string]CrawlerFunc {
	return registry
}

// RegisterNameFetcher registra uma função de busca de nome para um site específico
func RegisterNameFetcher(site string, fn NameFetcher) {
	nameRegistry[site] = fn
}

// FecthName busca o nome do produto de um site específico
func FecthName(site string, url string) (string, error) {
	if fn, ok := nameRegistry[site]; ok {
		return fn(url)
	}
	return "", nil
}
