package services

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/ViniMendes2515/price-crawler/internals/crawler"
	"github.com/ViniMendes2515/price-crawler/internals/historico"
	"github.com/ViniMendes2515/price-crawler/internals/notifier"
)

func agrupaSite(produtos []crawler.ProductInfo) map[string][]crawler.ProductInfo {
	grupos := make(map[string][]crawler.ProductInfo)

	for _, produto := range produtos {
		parsed, err := url.Parse(produto.URL)
		if err != nil {
			continue
		}
		dominio := parsed.Hostname()
		grupos[dominio] = append(grupos[dominio], produto)
	}

	return grupos
}

// Monitorar verifica se os produtos estão em promoção e envia notificações via Telegram
func Monitorar(produtos []crawler.ProductInfo, tg notifier.TelegramNotifier) {
	const telegramLimit = 4096
	const fatorDesvioPadrao = 0.6

	for dominio, produto := range produtos {
		mensagem := fmt.Sprintf("📢 Promoções em %s:\n\n", strconv.Itoa(dominio))

		promocao, err := historico.DetectarPromocao(produto.URL, produto.Price, fatorDesvioPadrao)
		if err != nil && !promocao {
			log.Println("Erro ao verificar promoção:", err)
			continue
		}

		bloco := fmt.Sprintf(
			"• %s\n💰 R$ %.2f\n🔗 %s\n\n",
			produto.Title,
			produto.Price,
			produto.URL,
		)

		if len(mensagem)+len(bloco) > telegramLimit {
			tg.Send(mensagem)
			mensagem = ""
		}

		mensagem += bloco

		historico.RegistrarPreco(produto.URL, produto.Price)

		if mensagem != fmt.Sprintf("📢 Promoções em %s:\n\n", strconv.Itoa(dominio)) {
			tg.Send(mensagem)
		}
	}

	if err := historico.Salvar(); err != nil {
		log.Println("Erro ao salvar histórico:", err)
	}
}
