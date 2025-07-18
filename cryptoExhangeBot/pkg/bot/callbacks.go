package bot

import (
	"fmt"
	"log"
	"strings"

	"cryptoChange/pkg/database"
	"cryptoChange/pkg/microservice"
	"cryptoChange/pkg/model"
	"cryptoChange/pkg/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleCallback обрабатывает входящие callback‑запросы
func HandleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	userID := callback.From.ID
	userState := state.GetState(userID)
	chatID := callback.Message.Chat.ID

	switch {
	case strings.HasPrefix(data, "coin:"):
		// Выбор криптовалюты через inline‑кнопки
		coin := strings.TrimPrefix(data, "coin:")
		if coin != "BTC" && coin != "ETH" && coin != "LTC" {
			return
		}
		userState.Coin = coin
		userState.Step = state.StepAmount
		text := "💵 Введите сумму оплаты:"
		markup := DefaultInlineKeyboard()
		EditOrSendMessage(bot, chatID, userState, text, markup)
		answerCallback(bot, callback, fmt.Sprintf("Вы выбрали %s", coin))
	case data == "confirm_order":
		processConfirmation(bot, callback, userState)
	case data == "main_menu":
		// Сброс состояния и возврат в главное меню
		state.DeleteState(userID)
		text := "🏠 Главное меню. Нажмите /start для начала нового заказа."
		editConfig := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, text)
		bot.Send(editConfig)
		answerCallback(bot, callback, "Главное меню")
	}
}

func answerCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, text string) {
	cfg := tgbotapi.NewCallback(callback.ID, text)
	bot.Request(cfg)
}

func processConfirmation(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, userState *state.UserState) {
	chatID := callback.Message.Chat.ID
	// Перед получением реквизита обращаемся к микросервису
	paymentDetail, err := microservice.GetPaymentDetail("http://your-microservice-url/api/get_payment_detail")
	if err != nil {
		paymentDetail = "Ошибка получения реквизитов ❌"
	}
	username := callback.From.UserName
	if username == "" {
		username = strings.TrimSpace(callback.From.FirstName + " " + callback.From.LastName)
	}
	order := &model.Order{
		Username:      username,
		Coin:          userState.Coin,
		PaymentAmount: userState.PaymentAmount,
		WalletAddress: userState.Wallet,
		PaymentDetail: paymentDetail,
	}
	if err = database.InsertOrder(order); err != nil {
		log.Printf("Ошибка при сохранении заявки: %v", err)
		text := "❌ Ошибка при сохранении заявки. Попробуйте позже."
		EditOrSendMessage(bot, chatID, userState, text, DefaultInlineKeyboard())
		state.DeleteState(callback.From.ID)
		answerCallback(bot, callback, "Ошибка")
		return
	}
	text := fmt.Sprintf("✅ Заявка создана!\n\n💰 Криптовалюта: %s\n💵 Сумма оплаты: %.2f руб\n📥 Кошелёк: %s\n🔑 Реквизит: %s",
		order.Coin, order.PaymentAmount, order.WalletAddress, order.PaymentDetail)
	markup := DefaultInlineKeyboard()
	EditOrSendMessage(bot, chatID, userState, text, markup)
	state.DeleteState(callback.From.ID)
	answerCallback(bot, callback, "Заказ подтвержден")
}
