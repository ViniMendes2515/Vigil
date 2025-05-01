package main

import (
	"fmt"
	"log"

	"github.com/ViniMendes2515/price-crawler/config"
	"github.com/ViniMendes2515/price-crawler/internals/crawler"
	"github.com/ViniMendes2515/price-crawler/internals/notifier"
	"github.com/ViniMendes2515/price-crawler/internals/services"
)

func main() {
	cfg := config.Load()

	tg, err := notifier.NewTelegramNotifier(cfg.TelegramToken, cfg.TelegramChatID)
	if err != nil {
		log.Fatal("Erro ao criar o bot do telegram: ", err)
		return
	}

	produtos, err := crawler.ScrapeKabum(cfg.KabumURLs)
	if err != nil {
		fmt.Println("Erro ao fazer o scraping: ", err)
		return
	}

	services.Monitorar(produtos, *tg)
}
