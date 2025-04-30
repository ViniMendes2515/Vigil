package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ViniMendes2515/price-crawler/config"
	"github.com/ViniMendes2515/price-crawler/internals/crawler"
	"github.com/ViniMendes2515/price-crawler/internals/notifier"
)

func main() {
	cfg := config.Load()

	for _, url := range cfg.KabumURLs {
		fmt.Println("🛒 Kabum URL:", url)
	}

	chatIdStr := os.Getenv("TELEGRAM_CHAT_ID")
	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil || chatId == 0 {
		log.Fatal("TELEGRAM_CHAT_ID nao foi denifido corretamente")
		return
	}

	tg, err := notifier.NewTelegramNotifier(cfg.TelegramToken, chatId)
	if err != nil {
		log.Fatal("Erro ao criar o bot do telegram: ", err)
		return
	}

	produtos, err := crawler.ScrapeKabum(cfg.KabumURLs)
	if err != nil {
		fmt.Println("Erro ao fazer o scraping: ", err)
		return
	}

	mensagemKabum := "📦 Produtos encontrados na Kabum:\n\n"

	for _, produto := range produtos {
		mensagemKabum += fmt.Sprintf("🛒 %s\n💰 R$ %.2f\n🔗 %s\n\n", produto.Title, produto.Price, produto.URL)
	}

	if err := tg.Send(mensagemKabum); err != nil {
		log.Fatal("Erro ao enviar mensagem para o Telegram: ", err)
		return
	}

}
