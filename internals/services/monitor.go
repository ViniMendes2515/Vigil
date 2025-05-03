package services

import (
	"fmt"
	"log"
	"net/url"

	"github.com/ViniMendes2515/price-crawler/internals/historico"
	"github.com/ViniMendes2515/price-crawler/internals/models"
	"github.com/ViniMendes2515/price-crawler/internals/notifier"
)

func agrupaSite(produtos []models.ProductInfo) map[string][]models.ProductInfo {
	grupos := make(map[string][]models.ProductInfo)

	for _, produto := range produtos {
		parsed, err := url.Parse(produto.Url)
		if err != nil {
			continue
		}
		dominio := parsed.Hostname()
		grupos[dominio] = append(grupos[dominio], produto)
	}

	return grupos
}

// Monitorar verifica se os produtos estão em promoção e envia notificações via Telegram
func Monitorar(produtos []models.ProductInfo, tg notifier.TelegramNotifier) {
	const telegramLimit = 4096
	const fatorDesvioPadrao = 0.6

	agrupados := agrupaSite(produtos)

	if err := historico.Carregar(); err != nil {
		log.Println("Erro ao carregar histórico:", err)
		return
	}

	for dominio, lista := range agrupados {
		mensagem := fmt.Sprintf("📢 Promoções em %s:\n\n", dominio)

		for _, produto := range lista {

			promocao, err := historico.DetectarPromocao(produto.Url, produto.Price, fatorDesvioPadrao)
			if err != nil && !promocao {
				log.Println("Erro ao verificar promoção:", err)
				continue
			}

			bloco := fmt.Sprintf(
				"• %s\n💰 R$ %.2f\n🔗 %s\n\n",
				produto.Title,
				produto.Price,
				produto.Url,
			)

			if len(mensagem)+len(bloco) > telegramLimit {
				tg.Send(mensagem)
				mensagem = ""
			}

			mensagem += bloco

			historico.RegistrarPreco(produto.Url, produto.Price)
		}

		if mensagem != "" {
			tg.Send(mensagem)
		}
	}

	if err := historico.Salvar(); err != nil {
		log.Println("Erro ao salvar histórico:", err)
	}
}
