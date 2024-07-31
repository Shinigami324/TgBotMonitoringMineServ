package main

import (
	"TgGraf/events/telegram"
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	var wg sync.WaitGroup
	onlineChan := make(chan int)

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
			msgArr := telegram.GetTextComamnd(update)

			if msgArr[0] == telegram.CommandOnline {
				wg.Add(1)

				go onlineServerFlow(update, bot, msgArr, &wg, onlineChan)
			}

			wg.Add(1)
			go telegram.DoCmd(update, bot, onlineChan)
		}
	}
}

func onlineServerFlow(update tgbotapi.Update, bot *tgbotapi.BotAPI, msgArr []string, wg *sync.WaitGroup, onlineChan chan int) {
	defer wg.Done()
	onlineServer, err := telegram.GetServerOnline(msgArr[1], update.Message.Chat.ID, bot)

	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Проблема в получение онлайна"))
	}

	onlineChan <- onlineServer
}
