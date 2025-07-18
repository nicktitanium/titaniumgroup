package logging

import (
	"fmt"
	"os"
	"runtime"
	"github.com/sirupsen/logrus"
)

const (
	logFileName     = "./.log"
	timeStampFormat = "2006-01-02 15:04:05"
)

var logFile *os.File

func InitLogger() {

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		logrus.Fatal()
	}

	logrus.SetOutput(logFile)

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: timeStampFormat,
	})

	logrus.SetLevel(logrus.DebugLevel)
}

func CloseLogFile() {
	logFile.Close()
}

func getCallerInfo() string {
	pc, file, line, ok := runtime.Caller(4)

	if !ok {
		return "unknown file"
	}

	function := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s %s:%d", file, function, line)
}

func HandlePanic() {
	if checkPanic := recover(); checkPanic != nil {
		fileAndLocation := getCallerInfo()
		logrus.Panic(fmt.Sprintf("panic: %v, location: %s", checkPanic, fileAndLocation))
	}
}
