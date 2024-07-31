package main

import (
	"TgGraf/events/telegram"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	fmt.Printf("Ответ от Redis: %s\n", pong)

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

				wg.Add(1)
				go bdRedis(rdb, msgArr[1], <-onlineChan, ctx, &wg)

				dataChannel := make(chan map[string]string)

				wg.Add(1)
				go func() {
					defer close(dataChannel)

					data, err := getRedisServerData(rdb, msgArr[1], ctx)

					if err != nil {
						fmt.Printf("Ошибка в извлечение данных: %s", err)
					}

					dataChannel <- data
				}()

				wg.Add(1)

				go func() {

					dataParse := <-dataChannel

					online, onlineExists := dataParse["online_players"]
					data, dataExists := dataParse["date"]

					if !onlineExists {
						fmt.Print("ключ online_players отсутствует в данных")
					}
					if !dataExists {
						fmt.Print("ключ date отсутствует в данных")
					}

					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Текущая дата мониторинга сервера: %s\n Текущий онлайн сервера %s", data, online)))
				}()
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

func bdRedis(rdb *redis.Client, serverName string, onlinePlayers int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	currentTime := time.Now().Format("02-01-2006")

	hashKey := serverName

	fieldValues := map[string]interface{}{
		"date":           currentTime,
		"online_players": onlinePlayers,
	}

	err := rdb.HMSet(ctx, hashKey, fieldValues).Err()

	if err != nil {
		log.Printf("Ошибка при добавление данных в Redis: %s", err)
	}
}

func getRedisServerData(rdb *redis.Client, serverName string, ctx context.Context) (map[string]string, error) {
	data, err := rdb.HGetAll(ctx, serverName).Result()

	if err != nil {
		return nil, fmt.Errorf("Ошибка извлечения данных: %s", err)
	}

	return data, nil
}
