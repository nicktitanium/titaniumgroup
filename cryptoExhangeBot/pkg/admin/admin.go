package admin

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Определите список ID администраторов (замените на свои)
var AdminIDs = map[int64]bool{
	6447535337: true, // замените 123456789 на реальный admin ID
}

// Rates хранит курсы для USDT и монет.
// Курс USDT – это, например, рубли за 1 USDT.
// Курсы монет – стоимость 1 монеты в USDT.
type Rates struct {
	USDT  float64
	Coins map[string]float64
}

var rates = Rates{
	USDT: 80.0, // по умолчанию: 80 руб/USDT
	Coins: map[string]float64{
		"BTC": 30000.0, // стоимость 1 BTC в USDT
		"ETH": 2000.0,
		"LTC": 100.0,
	},
}

// IsAdmin проверяет, является ли пользователь администратором.
func IsAdmin(userID int64) bool {
	return AdminIDs[userID]
}

// HandleAdminPanel выводит текущие курсы.
func HandleAdminPanel(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msgText := fmt.Sprintf("Текущие курсы:\nUSDT: %.2f руб/USDT\n", rates.USDT)
	for coin, rate := range rates.Coins {
		msgText += fmt.Sprintf("%s: %.2f USDT\n", coin, rate)
	}
	msgText += "\nИспользуйте:\n" +
		"/setusdt <значение> — для обновления курса USDT,\n" +
		"/setrate <coin> <значение> — для обновления курса монеты (например, /setrate BTC 35000).\n"
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	bot.Send(msg)
}

// HandleSetUSDT обновляет курс USDT.
func HandleSetUSDT(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := strings.Split(message.Text, " ")
	if len(args) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Использование: /setusdt <значение>")
		bot.Send(msg)
		return
	}
	rate, err := strconv.ParseFloat(args[1], 64)
	if err != nil || rate <= 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неверное значение курса USDT.")
		bot.Send(msg)
		return
	}
	rates.USDT = rate
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Курс USDT обновлён: %.2f руб/USDT", rate))
	bot.Send(msg)
}

// HandleSetRate обновляет курс конкретной монеты.
func HandleSetRate(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := strings.Split(message.Text, " ")
	if len(args) != 3 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Использование: /setrate <coin> <значение>")
		bot.Send(msg)
		return
	}
	coin := strings.ToUpper(args[1])
	rate, err := strconv.ParseFloat(args[2], 64)
	if err != nil || rate <= 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неверное значение курса для монеты.")
		bot.Send(msg)
		return
	}
	if _, ok := rates.Coins[coin]; !ok {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Монета %s не поддерживается.", coin))
		bot.Send(msg)
		return
	}
	rates.Coins[coin] = rate
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Курс для %s обновлён: %.2f USDT", coin, rate))
	bot.Send(msg)
}

// GetUSDT возвращает текущий курс USDT.
func GetUSDT() float64 {
	return rates.USDT
}
