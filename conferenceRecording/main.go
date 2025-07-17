package main

import (
	"conference_recognition/logging"
	webui "conference_recognition/web-ui"

	"github.com/sirupsen/logrus"
)

const logPath = "./.log"

func main() {

	logging.InitLogger(logPath)

	defer logging.CloseLogFile()
	defer logging.HandlePanic()
	defer logrus.Infof("server stoped")

	webui.StartWebUi()

}
