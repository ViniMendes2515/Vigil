package cmd

import (
	"fmt"
	"vigil/database"
	"vigil/internals/crawler"
	"vigil/internals/models"
	"vigil/internals/services"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Verifica se tem promoção nos produtos monitorados",
	Long:  `Verifica se tem promoção nos produtos monitorados e envia uma mensagem no Telegram.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("✅ Monitoramento iniciado com sucesso!")

		var produtos []models.ProductInfo

		for site, scrapeFunc := range crawler.GetRegistered() {
			urls, err := database.GetSiteUrls(site)
			if err != nil {
				fmt.Println("Erro ao obter URLs do banco de dados: ", err)
				continue
			}

			resultados, err := scrapeFunc(urls)
			if err != nil {
				fmt.Println("Erro ao fazer scraping: ", err)
				continue
			}

			produtos = append(produtos, resultados...)
		}

		tg, err := services.InitServices()
		if err != nil {
			fmt.Println("Erro ao inicializar serviços: ", err)
			return
		}

		services.Monitorar(produtos, *tg)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
