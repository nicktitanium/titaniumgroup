package tgbot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(message *tgbotapi.Message, wg *sync.WaitGroup) {

	defer wg.Done()

	if message.IsCommand() {
		var mt sync.Mutex

		switch message.Command() {

		case "start":
			//Воизбежании путаници будем временно блокировать доступ остальных рутин к обработки пользователя
			mt.Lock()
			logrus.Infof("user with id [%d] start bot", message.From.ID)
			b.HandleNewUser(message)
			mt.Unlock()

			b.sendWelcomeMessage(message.Chat.ID)
			fmt.Printf("[%v] - start Bot\n", message.From.UserName)
		case "help":
			b.sendTutorialMessage(message.Chat.ID, message.From.ID)

		default:
			b.sendUnknownCommandMessage(message.Chat.ID)
		}
		return
	}

	if b.Users[message.From.ID].State == awaitingPromptText {
		b.Users[message.From.ID].Prompt = message.Text
		b.deleteMessage(message.From.ID, message.MessageID)
		b.promptToSunoButtonMusicStyle(message.Chat.ID)
		b.Users[message.From.ID].State = awaitingSongStyle
		return
	} else if b.Users[message.From.ID].State == awaitingSongStyle {
		b.Users[message.From.ID].Prompt = fmt.Sprintf("%s %s", b.Users[message.From.ID].Prompt, message.Text)
		b.deleteMessage(message.From.ID, message.MessageID)
		b.isInstrumental(message.Chat.ID)
		b.Users[message.From.ID].State = awaitingInstumental
		return
	} else {
		b.sendMessage(message.Chat.ID, unknownMessage)
		return
	}

}
