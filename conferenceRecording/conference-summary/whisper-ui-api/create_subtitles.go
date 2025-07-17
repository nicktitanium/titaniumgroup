package whisperuiapi

import (
	"conference_recognition/config"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	startTransript = "/gradio_api/queue/join?"
	transcriptLogs = "/gradio_api/queue/data?session_hash="
)

type WhisperData struct {
	BaseUrl    string
	AudioURL   string
	Output     string
	AppData    string
	Parameters []interface{}
}

func getWhisperData() WhisperData {

	confData := config.GetConfigData("whisper_web_ui")

	return WhisperData{
		BaseUrl:    confData["url"].(map[interface{}]interface{})["base_url"].(string),
		AudioURL:   confData["url"].(map[interface{}]interface{})["audio_url"].(string),
		Output:     confData["path"].(map[interface{}]interface{})["output"].(string),
		AppData:    confData["path"].(map[interface{}]interface{})["app_data"].(string),
		Parameters: confData["parameters"].([]interface{}),
	}
}

func CreateSubtitles(audioPath string, temperature float32) (string, error) {

	confData := getWhisperData()

	client := http.Client{}

	audioFileInGradioDir, err := createNewEventDir(audioPath)

	if err != nil {
		return "", err
	}

	sessionID := uuid.New()

	fileName := strings.ReplaceAll(audioFileInGradioDir.Name(), "/", "\\")

	reqParams := RequestParameters{
		FilePath:    fileName,
		FileURL:     confData.AudioURL + fileName,
		OriginName:  filepath.Base(fileName),
		SessionHash: sessionID.String(),
		Temperature: temperature,
	}

	jsonPayload, err := reqParams.createPayload(confData.Parameters)

	if err != nil {
		return "", err
	}

	if err = sendAudioForTranscript(client, startTransript, jsonPayload); err != nil {
		return "", err
	}

	logrus.Debugf("info: wait transcript to finish ...")

	return getTranscript(client, transcriptLogs+sessionID.String(), filepath.Base(audioFileInGradioDir.Name()), reqParams.FileFormat)

}
