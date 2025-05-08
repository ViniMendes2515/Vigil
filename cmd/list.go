package cmd

import (
	"fmt"
	"vigil/database"

	"github.com/spf13/cobra"
)

var (
	prices bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista os produtos monitorados",
	Long:  `Lista todos os produtos monitorados com seus respectivos preços e URLs.`,
	Run: func(cmd *cobra.Command, args []string) {

		if prices {
			urlsPrice, err := database.ListPrices()
			if err != nil {
				fmt.Println("❌ Erro ao listar url com preços:", err)
				return
			}

			fmt.Println("✅ Produtos monitorados:")

			for _, row := range urlsPrice {
				fmt.Println(row)
			}

			return
		}

		urls, err := database.GetUrls()
		if err != nil {
			fmt.Println("❌ Erro ao listar URLs:", err)
			return
		}

		if len(urls) == 0 {
			fmt.Println("✅ Nenhum produto monitorado.")
			return
		}

		fmt.Println("✅ Produtos monitorados:")
		for url, id := range urls {
			fmt.Printf(" ID %d - URL : %s \n", id, url)
		}

	},
}

func init() {
	listCmd.Flags().BoolVarP(&prices, "price", "p", false, "Lista o preço de todos os produtos monitorados")
	rootCmd.AddCommand(listCmd)
}
