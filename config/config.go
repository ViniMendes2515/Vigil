package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken  string   // Token do bot do Telegram
	TelegramChatID int64    // ID do chat (usuário ou grupo) para envio
	KabumURLs      []string // URLs de produtos da Kabum
	AmazonURLs     []string // URLs de produtos da Amazon
	PriceLimit     float64  // Valor máximo para considerar um preço bom
}

func Load() Config {
	err := godotenv.Load("../.env")
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

	limitStr := os.Getenv("PRICE_LIMIT")
	limit, err := strconv.ParseFloat(limitStr, 64)
	if err != nil || limit <= 0 {
		limit = 1000 // Valor Default
	}

	// Carrega a lista de URLs da Kabum
	kabumStr := os.Getenv("KABUM_URLS")
	var kabum []string
	if kabumStr != "" {
		kabum = strings.Split(kabumStr, ",")
	}

	// Carrega a lista de URLs da Amazon
	amazonStr := os.Getenv("AMAZON_URLS")
	var amazon []string
	if amazonStr != "" {
		amazon = strings.Split(amazonStr, ",")
	}

	return Config{
		TelegramToken:  token,
		TelegramChatID: chatID,
		KabumURLs:      kabum,
		AmazonURLs:     amazon,
		PriceLimit:     limit,
	}

}
