package main

import (
	"log"
	"os"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Configuration struct {
	BotToken             string
	PsqlConnectionString string
}

type Reminder struct {
	gorm.Model
	ChatId   int64
	DateHour time.Time
	Content  string
}

func findHour(message string) string {
	var re = regexp.MustCompile(`(?m)\d{1,2}[h|:]\d{2}`)
	var result = re.FindString(message)
	return result
}

func findTomorrow(message string) bool {
	var re = regexp.MustCompile(`(?m)amanhã|amanha`)
	var result = re.MatchString(message)
	return result
}

func findToday(message string) bool {
	var re = regexp.MustCompile(`(?m)hoje`)
	var result = re.MatchString(message)
	return result
}

func main() {
	// Carregar as variáveis de ambiente
	var configuration Configuration
	configuration.BotToken = os.Getenv("TELEGRAM_TOKEN")
	configuration.PsqlConnectionString = os.Getenv("DATABASE_URL")

	// Conexão com o banco
	db, err := gorm.Open(postgres.Open(configuration.PsqlConnectionString), &gorm.Config{})

	// Realizar migração inicial
	db.AutoMigrate(&Reminder{})

	// Inicialização do bot
	bot, err := tgbotapi.NewBotAPI(configuration.BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

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
			var datetime_to_remember time.Time
			var content string
			chat_id := update.Message.Chat.ID
			datetime_now := time.Unix(int64(update.Message.Date), 0)
			response := tgbotapi.NewMessage(chat_id, "")

			switch update.Message.Command() {
			case "30m":
				response.Text = "Ok, te chamo daqui 30 minutos!"
				datetime_to_remember = datetime_now.Add(30 * time.Minute)
				content = "30 minutos já se passaram!"
			case "1h":
				response.Text = "Ok, te chamo daqui 1 hora!"
				datetime_to_remember = datetime_now.Add(30 * time.Hour)
				content = "1 hora já se passou!"
			case "1d":
				response.Text = "Ok, te chamo daqui 1 dia!"
				datetime_to_remember = datetime_now.AddDate(0, 0, 1)
				content = "1 dia já se passou!"
			}

			db.Create(&Reminder{ChatId: chat_id, Content: content, DateHour: datetime_to_remember})
			bot.Send(response)

		} else {
			// Salva a mensagem no banco
			chat_id := update.Message.Chat.ID
			content := update.Message.Text
			datetime := int64(update.Message.Date)
			unix_datetime := time.Unix(datetime, 0)
			db.Create(&Reminder{ChatId: chat_id, Content: content, DateHour: unix_datetime})

			// Responde a mensagem de volta pra quem enviou
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
