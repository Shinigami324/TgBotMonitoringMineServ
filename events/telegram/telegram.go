package telegram

import (
	grafcmd "TgGraf/events/grafCMD"
	"TgGraf/lib/er"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func DoCmd(update tgbotapi.Update, bot *tgbotapi.BotAPI) (err error) {
	defer func() { err = er.WrapIferr("Ошибка в выполнение команд: ", err) }()

	msgArr := strings.Fields(update.Message.Text)

	switch msgArr[0] {
	case commandStart:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, startBot))
	case commandHelp:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, helpBot))
	case commandNewGraf:
		grafCommand(update, bot)
	case commandOnline:

		online, err := getServerOnline(msgArr[1], update.Message.Chat.ID, bot)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Проблема в получение онлайна"))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Онлайн вашего сервера %s: %d", msgArr[1], online)))
	default:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Данная команда не известна"))
	}

	return nil
}

func grafCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) (err error) {
	defer func() { err = er.WrapIferr("Ошибка в выполнение команды /graf: ", err) }()

	if err := grafcmd.NewgraphInHours(); err != nil {
		log.Print("Ошибка в создание графика", err)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при создание графика"))
		return err
	}

	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, grafcmd.GetGraf())

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	if err != nil {
		log.Print("Ошибка при отправке графика", err)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при отправке графика"))
	}
	return nil
}

func getServerOnline(ip string, chatID int64, bot *tgbotapi.BotAPI) (online int, err error) {
	resp, err := http.Get(fmt.Sprintf("https://api.mcsrvstat.us/3/%s", ip))

	if err != nil {
		return 0, er.Wrap("Ошибка в отправке запроса: ", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("не удалось отправить запрос к htttp: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var serverData ServerData
	err = json.NewDecoder(resp.Body).Decode(&serverData)

	if err != nil {
		return 0, er.Wrap("Ошибка в парсинге Json: ", err)
	}

	if !serverData.Online {
		bot.Send(tgbotapi.NewMessage(chatID, "Сервер сейчас офлайн"))
		return 0, nil
	}

	return serverData.Player.Online, nil
}
