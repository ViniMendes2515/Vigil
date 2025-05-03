package main

import (
	"log"

	"github.com/ViniMendes2515/price-crawler/config"
	"github.com/ViniMendes2515/price-crawler/internals/notifier"
	"github.com/ViniMendes2515/price-crawler/internals/scheduler"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.Load()

	tg, err := notifier.NewTelegramNotifier(cfg.TelegramToken, cfg.TelegramChatID)
	if err != nil {
		log.Fatal("Erro ao criar o bot do telegram: ", err)
		return
	}

	c := cron.New()
	scheduler.BuscaPrecos(c, cfg, tg)
	c.Start()

	select {}
}
