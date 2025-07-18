package tgbot

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	usersDir     = "./users"
	userSongsDir = "songs"
	userClipsDir = "clips"
)

func (b *Bot) IsUserExist(userID int64) bool {
	_, exists := b.Users[userID]
	return exists
}

func (b *Bot) HandleNewUser(message *tgbotapi.Message) {
	if !b.IsUserExist(message.From.ID) {

		logrus.Infof("user with id [%d] not exist", message.From.ID)
		logrus.Debugf("start create new user (user id: %d)", message.From.ID)

		b.Users[message.From.ID] = &BotUser{
			State:       "",
			UserDirPath: fmt.Sprintf("%s%d", usersDir, message.From.ID),
		}

		logrus.Debugf("user struct was created (user id: %d)", message.From.ID)

		if dirName, err := b.CreateUserDir(message.From.ID); err != nil {
			logrus.Errorf("ERROR creating dir by path %s for new user (user id: %d), error_text: %v", dirName, message.From.ID, err)
			return
		}
	}
}

func (b *Bot) CreateUserDir(userID int64) (string, error) {

	err := os.Mkdir(b.Users[userID].UserDirPath, 0755)

	if err != nil {
		fmt.Println("Error create user's dir: ", err)
		return b.Users[userID].UserDirPath, err
	}

	err = os.Mkdir(fmt.Sprintf("%s/%s", b.Users[userID].UserDirPath, userSongsDir), 0755)

	if err != nil {
		fmt.Println("Error create user's song dir: ", err)
		return fmt.Sprintf("%s/%s", b.Users[userID].UserDirPath, userSongsDir), err
	}

	err = os.Mkdir(fmt.Sprintf("%s/%s", b.Users[userID].UserDirPath, userClipsDir), 0755)

	if err != nil {
		fmt.Println("Error create user's clips dir: ", err)
		return fmt.Sprintf("%s/%s", b.Users[userID].UserDirPath, userClipsDir), err
	}

	return "", nil
}
