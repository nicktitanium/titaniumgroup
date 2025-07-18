package bot

import (
	"cryptoChange/pkg/admin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleCommand обрабатывает входящие команды
func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		handleStart(bot, message)
	case "admin":
		if admin.IsAdmin(message.From.ID) {
			admin.HandleAdminPanel(bot, message)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "🚫 У вас нет доступа к админ-панели.")
			bot.Send(msg)
		}
	case "setusdt":
		if admin.IsAdmin(message.From.ID) {
			admin.HandleSetUSDT(bot, message)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "🚫 У вас нет доступа к этой команде.")
			bot.Send(msg)
		}
	case "setrate":
		if admin.IsAdmin(message.From.ID) {
			admin.HandleSetRate(bot, message)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "🚫 У вас нет доступа к этой команде.")
			bot.Send(msg)
		}
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "❓ Неизвестная команда. Попробуйте /start.")
		bot.Send(msg)
	}
}
