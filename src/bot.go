package main

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

type Configuration struct {
	BotToken             string
	PsqlConnectionString string
}

type Reminder struct {
	gorm.Model
	ChatId   string
	DateHour time.Time
	Content  string
}

func main() {
	// Carregar as variáveis de ambiente
	var configuration Configuration
	configuration.BotToken = os.Getenv("TELEGRAM_TOKEN")
	configuration.PsqlConnectionString = os.Getenv("DATABASE_URL")

	// Inicialização do bot
	bot, err := tgbotapi.NewBotAPI(configuration.BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Receber os updates do bot
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// Ignora updates que não sejam mensagens
		if update.Message == nil {
			continue
		}

		// Verifica se a mensagem é um comando:
		if update.Message.IsCommand() {
			chat_id := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(chat_id, "")

			switch update.Message.Command() {
			case "30m":
				msg.Text = "Ok, te chamo daqui 30 minutos!"
				// TODO: Salvar lembrete no banco
			case "1h":
				msg.Text = "Ok, te chamo daqui 1 hora!"
				// TODO: Salvar lembrete no banco
			case "1d":
				msg.Text = "Ok, te chamo daqui 1 dia!"
				// TODO: Salvar lembrete no banco
			}

			bot.Send(msg)
		}

		// Responde a mensagem de volta pra quem enviou
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

	}
}
