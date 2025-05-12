package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

var (
	deleteID uint
	all      bool
)

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Remove um produto do monitoramento",
	Long:  `Remove um produto do monitoramento, informando a URL.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if all {
			if err := repoUrls.RemoveAllUrls(ctx); err != nil {
				cmd.PrintErrln("Erro ao remover todos os produtos:", err)
				return
			}

			fmt.Println("✅ Todos os produtos removidos com sucesso.")
			return
		}

		if deleteID != 0 {
			if err := repoUrls.RemoveUrlById(ctx, deleteID); err != nil {
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

		urlStr := args[0]

		parsed, err := url.Parse(urlStr)
		if err != nil || parsed.Hostname() == "" {
			fmt.Println("❌ URL inválida")
			return
		}

		if err := repoUrls.RemoveUrl(ctx, urlStr); err != nil {
			cmd.PrintErrln("Erro ao remover URL:", err)
			return
		}

		fmt.Println("✅ URL removida com sucesso: ", urlStr)
	},
}

func init() {
	delCmd.Flags().UintVarP(&deleteID, "id", "i", 0, "ID do produto a ser removido")
	delCmd.Flags().BoolVarP(&all, "all", "a", false, "Remover todos os produtos")
	rootCmd.AddCommand(delCmd)
}
