package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"cryptoChange/pkg/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleStart инициирует диалог: генерирует капчу и отправляет приветственное сообщение
func handleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	// Генерируем капчу: два случайных числа от 1 до 10
	a := rand.Intn(10) + 1
	b := rand.Intn(10) + 1
	captchaAnswer := a + b

	// Инициализируем состояние пользователя
	userState := &state.UserState{
		Step:          state.StepCaptcha,
		CaptchaAnswer: captchaAnswer,
		ChatID:        chatID,
	}
	state.SetState(message.From.ID, userState)

	text := fmt.Sprintf("Привет! 😊 Добро пожаловать в крипто-обменник 🚀\n\n🔢 Решите капчу: %d + %d = ?", a, b)
	markup := DefaultInlineKeyboard()
	EditOrSendMessage(bot, chatID, userState, text, markup)
}

// handleState направляет обработку текстового сообщения в зависимости от текущего шага
func handleState(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userState := state.GetState(message.From.ID)
	switch userState.Step {
	case state.StepCaptcha:
		handleCaptcha(bot, message, userState)
	case state.StepCoinChoice:
		// На этом этапе предпочтительнее выбирать монету через inline‑кнопки,
		// но если пришёл текст, обрабатываем как fallback.
		handleCoinChoice(bot, message, userState)
	case state.StepAmount:
		handleAmount(bot, message, userState)
	case state.StepWallet:
		handleWallet(bot, message, userState)
	default:
		text := "❗️ Неизвестное состояние. Начните заново с /start"
		EditOrSendMessage(bot, userState.ChatID, userState, text, DefaultInlineKeyboard())
		state.DeleteState(message.From.ID)
	}
}

func handleCaptcha(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *state.UserState) {
	chatID := message.Chat.ID
	answer, err := strconv.Atoi(strings.TrimSpace(message.Text))
	if err != nil {
		text := "⚠️ Пожалуйста, введите число."
		EditOrSendMessage(bot, chatID, userState, text, DefaultInlineKeyboard())
		return
	}

	if answer != userState.CaptchaAnswer {
		text := "❌ Неправильный ответ. Попробуйте ещё раз."
		EditOrSendMessage(bot, chatID, userState, text, DefaultInlineKeyboard())
		return
	}

	// Переход к выбору криптовалюты
	userState.Step = state.StepCoinChoice
	text := "✅ Капча пройдена!\nВыберите криптовалюту: (BTC, ETH, LTC) 💰"
	markup := CoinSelectionKeyboard()
	EditOrSendMessage(bot, chatID, userState, text, markup)
}

func handleCoinChoice(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *state.UserState) {
	chatID := message.Chat.ID
	coin := strings.ToUpper(strings.TrimSpace(message.Text))
	if coin != "BTC" && coin != "ETH" && coin != "LTC" {
		text := "⚠️ Пожалуйста, выберите одну из предложенных опций: BTC, ETH, LTC"
		EditOrSendMessage(bot, chatID, userState, text, CoinSelectionKeyboard())
		return
	}
	userState.Coin = coin
	userState.Step = state.StepAmount
	text := "💵 Введите сумму оплаты:"
	EditOrSendMessage(bot, chatID, userState, text, DefaultInlineKeyboard())
}

func handleAmount(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *state.UserState) {
	chatID := message.Chat.ID
	amountStr := strings.TrimSpace(message.Text)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		text := "⚠️ Введите корректное положительное число для суммы."
		EditOrSendMessage(bot, chatID, userState, text, DefaultInlineKeyboard())
		return
	}
	userState.PaymentAmount = amount
	userState.Step = state.StepWallet
	text := "💵 Введите адрес кошелька, куда вы хотите получить монеты:"
	EditOrSendMessage(bot, chatID, userState, text, DefaultInlineKeyboard())
}

func handleWallet(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *state.UserState) {
	chatID := message.Chat.ID
	walletAddress := strings.TrimSpace(message.Text)
	userState.Wallet = walletAddress
	// Переходим к этапу подтверждения
	userState.Step = state.StepConfirmation
	text := fmt.Sprintf("Проверьте заявку:\n\n💰 Криптовалюта: %s\n💵 Сумма оплаты: %.2f руб\n📥 Кошелёк: %s\n\nВсе ли верно?",
		userState.Coin, userState.PaymentAmount, userState.Wallet)
	markup := ConfirmationKeyboard()
	EditOrSendMessage(bot, chatID, userState, text, markup)
}
