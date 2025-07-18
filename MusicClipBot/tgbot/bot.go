package tgbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"sai_project/logging"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	starttingOffset  = 0
	getUpdatesTimout = 60
	configFilePath   = "./config.yaml"
)

func GetConfigData(key string) (map[interface{}]interface{}, error) {

	configFileData, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		logrus.Fatalf("error reading config file. error_text: %v", err)
		return nil, err
	}

	var configData map[string]interface{}

	err = yaml.Unmarshal(configFileData, &configData)

	if err != nil {
		logrus.Fatalf("error unmarshaling config file data. error_text: %v", err)
		fmt.Println("Error unmarshaling config data")
		return nil, err
	}

	return configData[key].(map[interface{}]interface{}), err
}

func NewBot(token string) (*Bot, error) {

	api, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		logrus.Fatalf("error making new bot api. error_text: %v", err)
		return nil, err
	}

	return &Bot{
		Api:   api,
		Users: make(map[int64]*BotUser),
	}, nil
}

func (b *Bot) Start() {

	defer logging.HandlePanic()
	defer logging.CloseLogFile()

	log.Println("Bot started")

	update := tgbotapi.NewUpdate(starttingOffset)
	update.Timeout = getUpdatesTimout

	updates := b.Api.GetUpdatesChan(update)

	var wg sync.WaitGroup

	for update := range updates {

		if update.Message != nil {

			wg.Add(1)

			go func(message *tgbotapi.Message) {

				defer logging.HandlePanic()
				defer logging.CloseLogFile()

				b.handleMessage(message, &wg)

			}(update.Message)

		} else if update.CallbackQuery != nil {

			wg.Add(1)

			go func(callback *tgbotapi.CallbackQuery) {

				defer logging.HandlePanic()
				defer logging.CloseLogFile()

				b.handleCallback(callback, &wg)

			}(update.CallbackQuery)
		}
	}

	wg.Wait()
}
