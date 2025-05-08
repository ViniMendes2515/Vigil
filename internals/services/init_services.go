package services

import (
	"log"
	"vigil/config"
	"vigil/internals/notifier"
)

func InitServices() (*notifier.TelegramNotifier, error) {
	cfg := config.Load()

	tg, err := notifier.NewTelegramNotifier(cfg.TelegramToken, cfg.TelegramChatID)
	if err != nil {
		log.Fatal("Erro ao criar o bot do telegram: ", err)
		return nil, err
	}

	return tg, nil
}
