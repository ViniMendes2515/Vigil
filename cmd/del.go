package cmd

import (
	"fmt"
	"vigil/database"

	"github.com/spf13/cobra"
)

var (
	deleteID int
	all      bool
)

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Remove um produto do monitoramento",
	Long:  `Remove um produto do monitoramento, informando a URL.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			if err := database.RemoveAllUrls(); err != nil {
				cmd.PrintErrln("Erro ao remover todos os produtos:", err)
				return
			}

			fmt.Println("✅ Todos os produtos removidos com sucesso.")
			return
		}

		if deleteID != 0 {
			if err := database.RemoveUrlById(deleteID); err != nil {
				cmd.PrintErrln("Erro ao remover produto:", err)
				return
			}

			fmt.Printf("✅ Produto de ID %d removido com sucesso", deleteID)
			return
		}

		if len(args) == 0 {
			cmd.PrintErrln("Informe a URL ou use --id ou --all.")
			return
		}

		url := args[0]

		if err := database.RemoveUrl(url); err != nil {
			cmd.PrintErrln("Erro ao remover URL:", err)
			return
		}

		fmt.Println("✅ URL removida com sucesso: ", url)
	},
}

func init() {
	delCmd.Flags().IntVarP(&deleteID, "id", "i", 0, "ID do produto a ser removido")
	delCmd.Flags().BoolVarP(&all, "all", "a", false, "Remover todos os produtos")
	rootCmd.AddCommand(delCmd)
}
