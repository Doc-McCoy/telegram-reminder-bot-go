package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Configuration struct {
	BotApi                string `json:"bot_token"`
	MongoConnectionString string
}

// Interface do módulo que salva informações no mongoDB
type MongoData interface {
	Save() string
}

func recognizer(message string) {

}

func main() {
	// Carregar as variáveis do arquivo config.json
	var configuration Configuration
	byteFile, _ := ioutil.ReadFile("./config.dev.json")
	json.Unmarshal(byteFile, &configuration)
	configuration.MongoConnectionString = os.Getenv("MONGO_URL")

	// Inicialização do bot
	bot, err := tgbotapi.NewBotAPI(configuration.BotApi)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		// Ignora updates que não sejam mensagens
		if update.Message == nil {
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Verifica se a mensagem é um comando:
		if update.Message.IsCommand() {

			chat_id := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(chat_id, "")

			switch update.Message.Command() {
			case "30m":
				msg.Text = "Te chamo daqui 30 minutos!"
			case "1h":
				msg.Text = "Te chamo daqui 1 hora!"
			case "1d":
				msg.Text = "Te chamo daqui 1 dia!"
			}

			bot.Send(msg)
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)
	}
}
