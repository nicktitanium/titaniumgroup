package bot

import (
	"cryptoChange/pkg/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// EditOrSendMessage редактирует предыдущее сообщение, если оно уже было отправлено, или отправляет новое,
// при этом сохраняет MessageID в состоянии пользователя.
func EditOrSendMessage(bot *tgbotapi.BotAPI, chatID int64, userState *state.UserState, text string, markup tgbotapi.InlineKeyboardMarkup) (int, error) {
	// Если в состоянии уже хранится ID сообщения, пытаемся отредактировать его
	if userState.LastMessageID != 0 {
		editConfig := tgbotapi.NewEditMessageText(chatID, userState.LastMessageID, text)
		editConfig.ReplyMarkup = &markup
		_, err := bot.Send(editConfig)
		if err == nil {
			return userState.LastMessageID, nil
		}
	}
	// Если редактирование не удалось или сообщения ещё нет, отправляем новое
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup
	sentMsg, err := bot.Send(msg)
	if err != nil {
		return 0, err
	}
	userState.LastMessageID = sentMsg.MessageID
	return sentMsg.MessageID, nil
}

// MainMenuKeyboard возвращает inline‑клавиатуру с кнопкой «🏠 Главное меню»
func MainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	btn := tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu")
	row := tgbotapi.NewInlineKeyboardRow(btn)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

// CoinSelectionKeyboard возвращает inline‑клавиатуру для выбора криптовалюты плюс кнопку «🏠 Главное меню»
func CoinSelectionKeyboard() tgbotapi.InlineKeyboardMarkup {
	btnBTC := tgbotapi.NewInlineKeyboardButtonData("BTC", "coin:BTC")
	btnETH := tgbotapi.NewInlineKeyboardButtonData("ETH", "coin:ETH")
	btnLTC := tgbotapi.NewInlineKeyboardButtonData("LTC", "coin:LTC")
	mainMenu := tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu")
	row1 := tgbotapi.NewInlineKeyboardRow(btnBTC, btnETH, btnLTC)
	row2 := tgbotapi.NewInlineKeyboardRow(mainMenu)
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2)
}

// ConfirmationKeyboard возвращает клавиатуру для подтверждения заявки: «✅ Всё ок» и «🏠 Главное меню»
func ConfirmationKeyboard() tgbotapi.InlineKeyboardMarkup {
	btnConfirm := tgbotapi.NewInlineKeyboardButtonData("✅ Всё ок", "confirm_order")
	mainMenu := tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu")
	row := tgbotapi.NewInlineKeyboardRow(btnConfirm, mainMenu)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

// DefaultInlineKeyboard возвращает базовую клавиатуру с кнопкой «🏠 Главное меню»
func DefaultInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return MainMenuKeyboard()
}
