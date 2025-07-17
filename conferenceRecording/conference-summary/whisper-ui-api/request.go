package whisperuiapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"github.com/sirupsen/logrus"
)

const (
	baseContentType = "application/json"
)

type RequestParameters struct {
	OriginName  string
	FilePath    string
	FileURL     string
	SessionHash string
	FileFormat  string
	Temperature float32
}

const defaultTemp = 0.6

func (rq *RequestParameters) createPayload(parameters []interface{}) ([]byte, error) {

	payload := map[string]interface{}{
		"data": []interface{}{
			[]map[string]interface{}{
				{
					"meta": map[string]string{
						"_type": "gradio.FileData",
					},
					"mime_type": "audio/mpeg",
					"orig_name": filepath.Base(rq.FilePath),
					"path":      rq.FilePath,
					"url":       rq.FileURL,
				},
			},
		},
		"event_data":   nil,
		"fn_index":     0,
		"session_hash": rq.SessionHash,
		"trigger_id":   79,
	}

	for i, item := range parameters {

		if i == 3 {
			rq.FileFormat = item.(string)
		}

		if i == 17 && rq.Temperature != item {
			payload["data"] = append(payload["data"].([]interface{}), rq.Temperature)
			continue
		}

		payload["data"] = append(payload["data"].([]interface{}), item)

	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		logrus.Errorf("err_info: can't marshal request payload to json for Whisper-UI, err_text: %v", err)
		return nil, err
	}

	fmt.Println(string(jsonPayload))

	return jsonPayload, nil

}

func sendAudioForTranscript(client http.Client, endPoint string, payload []byte) error {

	confData := getWhisperData()

	logrus.Debugf("info: send request to Whisper-UI for transcript, url: %s/%s, payload: %s", confData.BaseUrl, endPoint, string(payload))

	resp, err := client.Post(confData.BaseUrl+endPoint, baseContentType, bytes.NewBuffer(payload))

	if err != nil {
		logrus.Fatalf("err_info: can't do POST request to %s%s, err_text: %v", confData.BaseUrl, endPoint, err)
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logrus.Errorf("err_info: can't read response body from  %s, err_text: %v", "http://127.0.0.1:7860/gradio_api/info", err)
		return err
	}

	if resp.StatusCode != 200 {
		logrus.Errorf("err_info: error status code from Whisper-UI: %d, response_body: %s", resp.StatusCode, string(respBody))
		return fmt.Errorf("Error status code from Whisper-UI: %d", resp.StatusCode)
	}

	logrus.Debugf("from_url: %s, response_body: %s", confData.BaseUrl+endPoint, string(respBody))

	return nil
}

//confData := config.GetConfigData("whisper_web_ui")
//confPayload := confData["parameters"].(map[interface{}]interface{})["data"].([]interface{})
//
//payload := map[string]interface{}{
//	"data": []interface{}{
//		[]map[string]interface{}{
//			{
//				"meta": map[string]string{
//					"_type": "gradio.FileData",
//				},
//				"mime_type": "au	dio/mpeg",
//				"orig_name": wd.OriginName,
//				"path":      wd.FilePath,
//				"url":       wd.FileURL,
//			},
//		},
//	},
//	"event_data":   nil,
//	"fn_index":     0,
//	"session_hash": wd.SessionHash,
//	"trigger_id":   0,
//}
//
//fmt.Println("Payloadывфы:", payload)
//
//for i := 0; i < len(confPayload); i++ {
//	if i == 4 {
//		wd.FileFormat = confPayload[i].(string)
//	}
//
//	if i == 18 && wd.Temperature != 0.0 {
//		confPayload[i] = wd.Temperature
//	} else if wd.Temperature == 0.0 {
//		wd.Temperature = defaultTemp
//	}
//
//	if val, ok := confPayload[i].(interface{}); ok {
//		payload["data"] = append(payload["data"].([]interface{}), val)
//	} else {
//		logrus.Warnf("Invalid type for confPayload at index %d", i)
//	}
//}
//
//fmt.Println("shit:", payload)
//
//jsonPayload, err := json.Marshal(payload)
//if err != nil {
//	logrus.Errorf("err_info: can't marshal request payload to json for Whisper-UI, err_text: %v", err)
//	return nil, err
//}
//
//fmt.Println("Payload:", payload)
//
//return jsonPayload, nil
