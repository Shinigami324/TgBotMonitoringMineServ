package main

import (
	"TgGraf/events/telegram"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6715777256:AAFujMk1cReHEm8gyPfHrSMFRNJTLyKluq8")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		fmt.Print("GetUpdatesChan errors: ", err)
	}

	for update := range updates {
		if update.Message != nil {
			telegram.DoCmd(update, bot)
		} else {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите тексотовое сообщение"))
		}

	}
}
