package tgbot

import (
	"fmt"
	"sync"

	//music "sai_project/suno_ai"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery, wg *sync.WaitGroup) {

	defer wg.Done()

	switch callback.Data {

	case startMakingClip:
		//b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.sendMessage(callback.Message.Chat.ID, promptMessage)
		b.Users[callback.From.ID].State = awaitingPromptText

	case rockStyle, popStyle, hiphopStyle, classicStyle:
		//b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.Users[callback.From.ID].Prompt = fmt.Sprintf("%s в стиле %s", b.Users[callback.From.ID].Prompt, callback.Data)
		b.Users[callback.From.ID].State = awaitingInstumental
		b.isInstrumental(callback.Message.Chat.ID)

	case instrumentalStyle:
		//b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.Users[callback.From.ID].IsInstrumental = true
		b.promptApproval(callback.Message.Chat.ID, callback.From.ID)

	case notinstrumentalStyle:
		//b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.Users[callback.From.ID].IsInstrumental = false
		b.promptApproval(callback.Message.Chat.ID, callback.From.ID)
		b.Users[callback.From.ID].State = awaitingPromptText

	case approvedPrompt:

		fmt.Printf("[%v] - %v\n", callback.From.UserName, b.Users[callback.From.ID].Prompt)
		b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.sendMessage(callback.Message.Chat.ID, waitSongMessage)

		// err := b.SongGeneration(callback.Message.Chat.ID, callback.From.ID, b.Users[callback.From.ID].Prompt, b.Users[callback.From.ID].IsInstrumental)

		// if err != nil {
		// 	fmt.Printf("Error: %v\n\n", err)
		// 	b.badRequest(callback.Message.Chat.ID)
		// 	return
		// }
		//b.AwaitingSongFiles = false

		b.chooseSong(callback.Message.Chat.ID)
		b.Users[callback.From.ID].State = awaitingVideoModel

		// Текст песни для тестов
		b.Users[callback.From.ID].Lyrics = "[Verse]\nСнова вместе собрались\nВ звуке старых добрых дней\nСтруны гитарных песен\nЗвучат в сердце у людей\n\n[Verse 2]\nВспомним юность без забот\nМузыка ведёт вперёд\nРуки вверх\nВоздух дрожит\nВечный рок нас вновь манит\n\n[Chorus]\nПусть играет наш мотив\nСтарый рок соединит\nЭти стены вновь поют\nСтарый друг\nНаш общий путь\n\n[Verse 3]\nГолос в микрофон летит\nКаждый звук - наш новый хит\nЛента закрутит момент\nПод гитарный инструмент\n\n[Chorus]\nПусть играет наш мотив\nСтарый рок соединит\nЭти стены вновь поют\nСтарый друг\nНаш общий путь\n\n[Bridge]\nВ каждом аккорде живёт\nЭхо прошлых нежных нот\nМузыка нас сохранит\nЭтот мир нас вдохновит"

	case notapprovedPrompt:
		b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.sendMessage(callback.Message.Chat.ID, "Пожалуйста, введите запрос заново.")
		b.Users[callback.From.ID].State = awaitingPromptText

	case firstSong:
		b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.isAnimated(callback.From.ID)
		b.Users[callback.From.ID].SongNumber = 1
		fmt.Printf("User [%s] choose song number: %d\n", callback.From.UserName, b.Users[callback.From.ID].SongNumber)

	case secondSong:
		b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.isAnimated(callback.From.ID)
		b.Users[callback.From.ID].SongNumber = 2
		fmt.Printf("User [%s] choose song number: %d\n", callback.From.UserName, b.Users[callback.From.ID].SongNumber)
		b.Users[callback.From.ID].State = awaitingVideoModel

	// case realisticVideo:

	// 	b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	// 	b.sendMessage(callback.Message.Chat.ID, "Принято! Начинаю генерацию видео :)\n Примерное время ожидания 30 минут...")

	// 	b.UserIsAnimated[callback.From.ID] = false

	// 	fmt.Print("[%s]Is animated video: ", callback.From.UserName, b.UserIsAnimated)

	// 	err := b.VideoGeneration(b.SongNumber, b.SongDuration[b.SongNumber-1], callback.Message.Chat.ID)

	// 	if err != nil {
	// 		fmt.Printf("Error: %v\n\n", err)
	// 		b.badRequest(callback.Message.Chat.ID)
	// 		return
	// 	}

	case animatedVideo:
		b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.sendMessage(callback.Message.Chat.ID, "Принято! Начинаю генерацию видео :)\n Примерное время ожидания 30 минут...")
		b.Users[callback.From.ID].isAnimated = true

		fmt.Printf("[%s]Is animated video: %v\n", callback.From.UserName, b.Users[callback.From.ID].isAnimated)

		err := b.CreateVideoClips(callback.From.ID) //b.SongDuration[b.SongNumber-1]

		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			b.badRequest(callback.Message.Chat.ID)
			return
		}

		b.Users[callback.From.ID].State = awaitingVideo

	case regenerateVideo:

		fmt.Println("Song number: ", b.Users[callback.From.ID].SongNumber, callback.Data)
		b.sendMessage(callback.Message.Chat.ID, "Принято! Начинаю генерацию видео :)\n Примерное время ожидания 30 минут...")

		err := b.CreateVideoClips(callback.Message.Chat.ID)

		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			b.badRequest(callback.Message.Chat.ID)
			return
		}
		b.Users[callback.From.ID].State = awaitingVideo

	case stopGeneraton:
		b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.sendMessage(callback.Message.Chat.ID, endMessage)

	default:
		b.sendMessage(callback.Message.Chat.ID, unknownMessage)
	}
}
