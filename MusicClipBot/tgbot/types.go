package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Bot struct {
	Users map[int64]*BotUser
	Api   *tgbotapi.BotAPI
}

type BotUser struct {
	State          string
	UserDirPath    string
	Prompt         string
	IsInstrumental bool
	Lyrics         string
	SongDuration   []float64
	SongNumber     int
	isAnimated     bool
}
