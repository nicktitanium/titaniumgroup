package llmstudioapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func PostRequestToLLM(pageName string, payload map[string]interface{}) (map[string]interface{}, error) {

	confData := getLmStudioData()

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		logrus.Errorf("err_info: can't marshal payload for %s/%s, err_text: %v", confData.BaseURL, pageName, err)
		return nil, err
	}

	//logrus.Debugf("payload_to_llm: %s", string(jsonPayload))

	client := http.Client{}

	requestUrl := fmt.Sprintf("%s/%s", confData.BaseURL, pageName)

	resp, err := client.Post(requestUrl, "application/json", bytes.NewBuffer(jsonPayload))

	if err != nil {
		logrus.Fatalf("err_info: can't do post request, err_text: %v", err)
		return nil, err
	}

	logrus.Debug("start creating subtitles summary with LM Studio API")

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		logrus.Errorf("err_info: error reading response body from %s/%s, err_text: %v", confData.BaseURL, pageName, err)
		return nil, err
	}

	respJson := map[string]interface{}{}

	err = json.Unmarshal(respBody, &respJson)

	if err != nil {
		logrus.Errorf("err_info: can't unmarshaling response to JSON from %s/%s, err_text: %v", confData.BaseURL, pageName, err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		logrus.Errorf("err_info: error status code from LM Studio: %d, response_body: %s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("Error status code from LM Studio: %d", resp.StatusCode)
	}

	logrus.Infof("summary successfully created!")

	return respJson, nil
}
