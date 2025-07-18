package bot

import (
	"log"
	"math/rand"
	"time"

	"cryptoChange/pkg/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start инициализирует бота и запускает цикл обработки обновлений
func Start(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("🤖 Авторизовались как: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	rand.Seed(time.Now().UnixNano())

	for update := range updates {
		// Если пришёл callback‑запрос (нажатие inline‑кнопки)
		if update.CallbackQuery != nil {
			HandleCallback(bot, update.CallbackQuery)
			continue
		}
		// Если сообщение отсутствует – пропускаем
		if update.Message == nil {
			continue
		}
		// Если сообщение – команда, обрабатываем её
		if update.Message.IsCommand() {
			handleCommand(bot, update.Message)
			continue
		}
		// Если состояние не установлено – просим начать с /start
		if state.GetState(update.Message.From.ID) == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❗️ Пожалуйста, начните с команды /start")
			bot.Send(msg)
			continue
		}
		// Обработка обычных текстовых сообщений (в рамках диалога)
		handleState(bot, update.Message)
	}
}
