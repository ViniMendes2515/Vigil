package notifier

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	Bot    *tgbotapi.BotAPI
	ChatId int64
}

func NewTelegramNotifier(token string, chatId int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramNotifier{
		Bot:    bot,
		ChatId: chatId,
	}, nil

}

func (t *TelegramNotifier) Send(message string) error {
	msg := tgbotapi.NewMessage(t.ChatId, message)
	_, err := t.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
