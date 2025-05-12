package cmd

import (
	"vigil/database"

	"github.com/spf13/cobra"
)

var (
	repoUrls    database.UrlRepository
	repoHistory database.PriceHistoryRepository
)

var rootCmd = &cobra.Command{
	Use:   "vigil",
	Short: "👁 Vigil — seu vigia de promoções em Go",
	Long: `Vigil é um rastreador inteligente de preços desenvolvido em Go, 
	projetado para monitorar ofertas online e alertar automaticamente quando o menor valor histórico é atingido.
	O Vigil está sempre atento. 👁`,
}

// Execute executa o comando raiz
func Execute() error {
	return rootCmd.Execute()
}

// SetRepositories define os repositórios a serem usados
func SetRepositories(urlRepo database.UrlRepository, historyRepo database.PriceHistoryRepository) {
	repoUrls = urlRepo
	repoHistory = historyRepo
}
