package cmd

import (
	"context"
	"vigil/pkg/utils"

	"github.com/spf13/cobra"
)

var showID int

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Mostra informações de um produto especifico",
	Long:  `Mostra informações de um produto especifico, informando a URL como argumento`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) > 0 && showID == 0 {
			url = args[0]
		}

		if showID != 0 {
			detalhes, err := repoHistory.ShowByID(context.Background(), showID)
			if err != nil {
				cmd.PrintErrln("❌ Erro ao buscar informações do produto:", err)
				return
			}
			utils.PrintProductTable(*detalhes)
			return
		}

		detalhes, err := repoHistory.ShowByUrl(context.Background(), url)
		if err != nil {
			cmd.PrintErrln("❌ Erro ao buscar informações do produto:", err)
			return
		}

		utils.PrintProductTable(*detalhes)
	},
}

func init() {
	showCmd.Flags().IntVarP(&showID, "id", "i", 0, "ID do produto")
	rootCmd.AddCommand(showCmd)
}
