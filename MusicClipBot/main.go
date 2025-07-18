package main

import (
	"github.com/sirupsen/logrus"
	"sai_project/logging"
	"sai_project/tgbot"
)

func main() {

	defer logging.HandlePanic()

	logging.InitLogger()

	defer logging.CloseLogFile()

	tgConfData, err := tgbot.GetConfigData("telegram")

	if err != nil {
		return
	}

	bot, err := tgbot.NewBot(tgConfData["token"].(string))
	if err != nil {
		logrus.Fatalf("Error creating bot: %v", err)
	} else {
		logrus.Info("Bot was created")
	}

	bot.Start()

}
