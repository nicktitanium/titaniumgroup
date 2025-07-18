package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
)

func (b *Bot) sendWelcomeMessage(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, welcomeMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userStartMakingClip, startMakingClip),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) sendTutorialMessage(chatID, userID int64) {
	msg := tgbotapi.NewMessage(chatID, tutorialMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userStartMakingClip, startMakingClip),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) sendUnknownCommandMessage(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, unknownComandMessage)
	b.Api.Send(msg)
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.Api.Send(msg)
}

func (b *Bot) deleteMessage(chatID int64, messageID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := b.Api.Request(deleteMsg)
	if err != nil {
		logrus.Printf("Error deleting message: %v", err)
	}
}

func (b *Bot) promptToSunoButtonMusicStyle(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, songStyleMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userRockStyle, rockStyle),
			tgbotapi.NewInlineKeyboardButtonData(userHipHopStyle, hiphopStyle),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userPopStyle, popStyle),
			tgbotapi.NewInlineKeyboardButtonData(userClassicStyle, classicStyle),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) isInstrumental(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, instrumentalMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userNotInstrumentalStyle, notinstrumentalStyle),
			tgbotapi.NewInlineKeyboardButtonData(userInstrumentalStyle, instrumentalStyle),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) ContinueGenerate(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, instrumentalMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userNotInstrumentalStyle, notinstrumentalStyle),
			tgbotapi.NewInlineKeyboardButtonData(userInstrumentalStyle, instrumentalStyle),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) promptApproval(chatID, userID int64) {

	var isWithVoice string

	// Так как isInstrumental означает: будет ли песня исключительно инструментальной композицией или нет, где false - означет налчиие вокала -> наличие вокала - isInstrumental == false
	if !b.Users[userID].IsInstrumental {
		isWithVoice = "да"
	} else {
		isWithVoice = "нет"
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ваш запрос: %s\nС вокалом:  %s\nПодтвердите, пожалуйста.", b.Users[userID].Prompt, isWithVoice))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userApprovedPrompt, approvedPrompt),
			tgbotapi.NewInlineKeyboardButtonData(userNotApprovedPrompt, notapprovedPrompt),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) chooseSong(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, chooseSongMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userFirstSong, firstSong),
			tgbotapi.NewInlineKeyboardButtonData(userSecondSong, secondSong),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) badRequest(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, errorGenerateMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userRegenerateVideo, regenerateVideo),
			tgbotapi.NewInlineKeyboardButtonData(userStopGeneraton, stopGeneraton),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) ApprovalVideo(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, aproveVideo)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userApprovedVideo, approvedVideo),
			tgbotapi.NewInlineKeyboardButtonData(userNotApprovedVideo, notapprovedVideo),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) isAnimated(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, videoStyleMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userAnimatedVideo, animatedVideo),
			tgbotapi.NewInlineKeyboardButtonData(userRealisticVideo, realisticVideo),
		),
	)
	b.Api.Send(msg)
}

func (b *Bot) sendAudio(f *os.File, id int64) error {

	file, err := os.Open(f.Name())

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	r := tgbotapi.FileReader{Name: f.Name(), Reader: file}

	audio := tgbotapi.NewAudio(id, r)

	b.Api.Send(audio)

	return nil

}

func (b *Bot) sendVideo(filename string, id int64) error {
	// Открываем видеофайл
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening video file: %v\n", err)
		return err
	}
	defer file.Close()

	fmt.Println("Video file opened successfully")

	r := tgbotapi.FileReader{Name: file.Name(), Reader: file}

	// Создаем новый объект для отправки видео
	videoClip := tgbotapi.NewVideo(id, r)

	// Отправляем видео
	fmt.Println("Sending video file")
	_, err = b.Api.Send(videoClip)
	if err != nil {
		fmt.Printf("Error sending video: %v\n", err)
		return err
	}

	fmt.Println("Video sent successfully")

	fmt.Println("Removing video file")
	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("Error removing video file: %v\n", err)
	} else {
		fmt.Println("Video file removed successfully")
	}

	return nil
}
