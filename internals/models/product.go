package models

type ProductInfo struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
	Url   string  `json:"url"`
	Name  string  `json:"name"`
}

type HistoricoPreco struct {
	Preco float64 `json:"preco"`
	Data  string  `json:"data"`
}

type ProductDetails struct {
	Site              string           `json:"site"`
	Nome              string           `json:"nome"`
	Url               string           `json:"url"`
	PrecoLimite       float64          `json:"preco_limite"`
	PrecoAtual        float64          `json:"preco_atual"`
	PrecoMinimo       float64          `json:"preco_minimo"`
	PrecoMaximo       float64          `json:"preco_maximo"`
	PrecoMedio        float64          `json:"preco_medio"`
	UltimaVerificacao string           `json:"ultima_verificacao"`
	TotalColetas      int              `json:"total_coletas"`
	HistoricoRecentes []HistoricoPreco `json:"historico_recentes"`
}
