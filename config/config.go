package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken  string // Token do bot do Telegram
	TelegramChatID int64  // ID do chat (usuário ou grupo) para envio
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("❌ TELEGRAM_TOKEN não está definido.")
	}

	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("❌ TELEGRAM_CHAT_ID inválido ou ausente.")
	}

	return Config{
		TelegramToken:  token,
		TelegramChatID: chatID,
	}

}
