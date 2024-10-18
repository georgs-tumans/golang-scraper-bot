package bot_fixer

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotFixer) SendMessage(chatId int64, text string, entities []tgbotapi.MessageEntity) {
	msg := tgbotapi.NewMessage(chatId, text)
	if len(entities) > 0 {
		msg.Entities = entities
	}
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := b.Bot.Send(msg)

	if err != nil {
		log.Printf("[Bot fixer] Error sending a message: %s", err.Error())
	}

	log.Printf("[Bot fixer] Sent message to chat: %d.Message: %s", chatId, text)
}

func (b *BotFixer) SendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := b.Bot.Send(msg)

	return err
}
