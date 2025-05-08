package utils

import (
	"fmt"
	"os"
	"vigil/internals/models"

	"github.com/olekukonko/tablewriter"
)

func PrintProductTable(detalhes models.ProductDetails) {
	fmt.Printf("\n🔍 Informações do produto:\n\n")
	data := [][]string{
		{"🛒 Site", detalhes.Site},
		{"🔗 URL", detalhes.Url},
		{"🧾 Nome", detalhes.Nome},
		{"💰 Preço Atual", fmt.Sprintf("R$ %.2f", detalhes.PrecoAtual)},
		{"🏷 Preço Limite", fmt.Sprintf("R$ %.2f", detalhes.PrecoLimite)},
		{"📉 Menor Preço", fmt.Sprintf("R$ %.2f", detalhes.PrecoMinimo)},
		{"📈 Maior Preço", fmt.Sprintf("R$ %.2f", detalhes.PrecoMaximo)},
		{"📆 Última Verificação", detalhes.UltimaVerificacao},
		{"📊 Total de Coletas", fmt.Sprintf("%d", detalhes.TotalColetas)},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Campo", "Valor"})
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)

	table.SetAutoWrapText(false)
	table.SetColumnSeparator("│")
	table.SetCenterSeparator("┼")
	table.SetRowSeparator("─")
	table.SetHeaderLine(true)

	table.AppendBulk(data)
	table.Render()

	if len(detalhes.HistoricoRecentes) > 0 {
		fmt.Println("\n📈 Histórico de preços (últimos 5):")

		histTable := tablewriter.NewWriter(os.Stdout)
		histTable.SetHeader([]string{"Data", "Preço"})
		for _, h := range detalhes.HistoricoRecentes {
			histTable.Append([]string{h.Data, fmt.Sprintf("R$ %.2f", h.Preco)})
		}
		histTable.SetAlignment(tablewriter.ALIGN_LEFT)
		histTable.Render()
	}
}
