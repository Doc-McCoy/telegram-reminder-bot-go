package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
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

func findDate(message string) string {
	var re = regexp.MustCompile(`(?m)\d{2}\/\d{2}\/\d{4}`)
	var result = re.FindString(message)
	return result
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
			// Extrair informações da mensagem
			content_ok := true
			chat_id := update.Message.Chat.ID
			content := update.Message.Text
			response := tgbotapi.NewMessage(chat_id, "")
			var date time.Time
			// datetime := int64(update.Message.Date)
			// unix_datetime := time.Unix(datetime, 0)

			date_find := findDate(content)
			hour_find := findHour(content)
			tomorrou_find := findTomorrow(content)
			today_find := findToday(content)

			// Localiza e define data
			if date_find == "" {
				if today_find {
					date = time.Now()
				} else if tomorrou_find {
					now := time.Now()
					date = now.AddDate(0, 0, 1)
				} else {
					content_ok = false
				}
			}

			// Localiza e define hora
			if hour_find == "" {
				content_ok = false
			}

			if content_ok {
				split_time := strings.Split(hour_find, ":")
				date.Hour = split_time[0]
				date.Minute = split_time[1]

				fmt.Println(date)

				// Salva infos no banco
				// db.Create(&Reminder{ChatId: chat_id, Content: content, DateHour: unix_datetime})
				response.Text = "Lembrete salvo para dia X, as Y."
			} else {
				response.Text = "Desculpe, não consegui identificar a data e a hora do lembrete."
			}
			// Responde a mensagem de volta pra quem enviou
			bot.Send(response)
		}
	}
}
