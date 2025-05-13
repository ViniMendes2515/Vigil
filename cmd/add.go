package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"vigil/internals/crawler"

	"github.com/spf13/cobra"
)

var (
	precoLimite float64
	precoAtual  float64
	filePath    string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adiciona um produto para monitorar",
	Long:  `Adiciona um produto para monitorar, informando a URL, o preco inicial e o preço limite se quiser.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if filePath != "" {
			readFile(ctx, filePath)
			return
		}

		if len(args) != 1 {
			cmd.PrintErrln("Aviso: Informe a URL como argumento principal ou use --file seguido do caminho de um arquivo CSV com o padrão (url, preco, limite)")
			return
		}

		endereco := args[0]

		parsed, err := url.Parse(endereco)
		if err != nil || parsed.Hostname() == "" {
			fmt.Println("❌ URL inválida")
			return
		}

		site := strings.ToLower(parsed.Hostname())

		name, err := crawler.FecthName(site, endereco)
		if err != nil {
			fmt.Println("❌ Erro ao buscar nome do produto:", err)
			return
		}

		err = repoUrls.AddUrl(ctx, site, endereco, name, precoAtual, precoLimite)
		if err != nil {
			fmt.Println("❌ Erro ao adicionar URL:", err)
			return
		}

		fmt.Printf("✅ Produto adicionado com sucesso: %s (Preço: R$%.2f) (Limite: R$%.2f)\n", endereco, precoAtual, precoLimite)
	},
}

func init() {
	addCmd.Flags().Float64VarP(&precoLimite, "limite", "l", 0.0, "Preço limite para o produto")
	addCmd.Flags().Float64VarP(&precoAtual, "preco", "p", 0.0, "Preço inicial do produto")
	addCmd.Flags().StringVarP(&filePath, "file", "f", "", "Caminho para arquivo CSV com produtos. Padrão CSV (url, preco, limite)")

	rootCmd.AddCommand(addCmd)
}

// readFile lê um arquivo CSV e adiciona os produtos ao banco de dados
func readFile(ctx context.Context, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("❌ Erro ao abrir arquivo:", err)
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("❌ Erro ao ler arquivo:", err)
		return
	}

	if (len(records)) < 2 {
		fmt.Println("⚠️ Nenhuma linha encontrada no CSV.")
		return
	}

	for i, linha := range records[1:] {
		if (len(linha)) < 2 {
			fmt.Printf("⚠️ Linha %d incompleta. Ignorada.\n", i+2)
			continue
		}

		endereco := strings.TrimSpace(linha[0])
		precoStr := strings.TrimSpace(linha[1])
		limiteStr := ""
		if len(linha) >= 3 {
			limiteStr = strings.TrimSpace(linha[2])
		}

		preco, err := strconv.ParseFloat(precoStr, 64)
		if err != nil {
			fmt.Printf("⚠️ Linha %d: preço inválido: %s\n", i+2, precoStr)
			continue
		}

		limite := 0.0
		if limiteStr != "" {
			limite, _ = strconv.ParseFloat(limiteStr, 64)
		}

		parsed, err := url.Parse(endereco)
		if err != nil || parsed.Hostname() == "" {
			fmt.Printf("❌ Linha %d: URL inválida: %s\n", i+2, endereco)
			continue
		}

		site := strings.ToLower(parsed.Hostname())

		name, err := crawler.FecthName(site, endereco)
		if err != nil {
			fmt.Println("❌ Erro ao buscar nome do produto:", err)
			return
		}

		err = repoUrls.AddUrl(ctx, site, endereco, name, preco, limite)
		if err != nil {
			fmt.Println("❌ Erro ao adicionar URL:", err)
			continue
		}

		fmt.Printf("✅ Linha %d: adicionada %s (Preço: %.2f | Limite: %.2f)\n", i+2, endereco, preco, limite)
	}
}
