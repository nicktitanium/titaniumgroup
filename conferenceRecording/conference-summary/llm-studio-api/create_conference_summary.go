package llmstudioapi

import (
	"conference_recognition/config"
	"errors"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

const (
	complitionWithChat = "chat/completions"
)

type LLMData struct {
	BaseURL       string
	SystemMessage string
	UserMessage   string
}

func getLmStudioData() LLMData {

	confData := config.GetConfigData("lm_studio")

	return LLMData{
		BaseURL:       confData["url"].(map[interface{}]interface{})["base_url"].(string),
		SystemMessage: confData["msg"].(map[interface{}]interface{})["system"].(string),
		UserMessage:   confData["msg"].(map[interface{}]interface{})["user"].(string),
	}

}

func CreateSummary(transcriptFilePath, prompt string) (string, error) {

	confData := getLmStudioData()

	if prompt != "" {
		confData.UserMessage = prompt
	}

	subtitles, err := ioutil.ReadFile(transcriptFilePath)

	if err != nil {
		logrus.Errorf("err_info: can't find transcript file for summary, file_path: %s, err_text: %v", transcriptFilePath, err)
	}

	if subtitles == nil {
		logrus.Error("err_info: bad whisper request (transcript file is null)")
		return "", errors.New("bad whisper request (transcript file is null)")
	}

	payload := map[string]interface{}{
		"model": "granite-3.0-2b-instruct",
		"messages": []map[string]string{
			{"role": "system", "content": confData.SystemMessage},
			{"role": "user", "content": confData.UserMessage + "\n" + string(subtitles)},
		},
		"temperature": 0.7,
		"max_tokens":  32000,
		"stream":      false,
	}

	resp, err := PostRequestToLLM(complitionWithChat, payload)

	if err != nil {
		return "", err
	}

	choices := resp["choices"].([]interface{})

	firstChoice := choices[0].(map[string]interface{})

	message := firstChoice["message"].(map[string]interface{})

	content := message["content"].(string)

	return content, nil
}
