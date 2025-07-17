package whisperuiapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Format string

const (
	AVI Format = "avi"
	MP4 Format = "mp4"
	MOV Format = "MOV"
)

func getTranscript(client http.Client, endPoint, audioFileName, fileFormat string) (string, error) {

	confData := getWhisperData()

	if fileFormat == "SRT" {
		fileFormat = strings.ToLower(fileFormat)
	}

	resp, err := client.Get(confData.BaseUrl + endPoint)

	if err != nil {
		logrus.Errorf("err_info: can't do GET request to %s%s, err_text: %v", confData.BaseUrl, endPoint, err)
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logrus.Errorf("err_info: can't read response body from  %s, err_text: %v", confData.BaseUrl+endPoint, err)
		return "", err
	}

	if resp.StatusCode != 200 {
		logrus.Errorf("err_info: can't read response body from  %s/%s, response_body: %s, err_text: %v", confData.BaseUrl, "gradio_api/info", string(respBody), err)
		return "", err
	}

	//transcriptedFilePath := fmt.Sprintf("%s/%s.%s", confData.Output, strings.TrimSuffix(audioFileName, ".mp4"), fileFormat)

	newFilePath, err := TrimFileFormat(audioFileName)

	if err != nil {
		return "", err
	}

	transcriptedFilePath := fmt.Sprintf("%s/%s%s", confData.Output, newFilePath, fileFormat)

	logrus.Debugf("file successfully transcripted. go to  %s", transcriptedFilePath)

	return transcriptedFilePath, nil

}

//C:\Users\sayap\AppData\Local\Temp\gradio\efe751fc684424bf01cb9893ed21017a682904d12cfaceaeee06ae5ac9bcfe978c33be2f12057a0dd34fdab608a5d6d5fe21d917ae37d3055481b2aef4703039
//C:\Users\sayap\AppData\Local\Temp\gradio\80724b9308a27df0333eb1d75bf35ba9a5eb8ce4886ed2122cc75c1f89bf149416214aecf89104465a41784071c0e91e405930c1025a53ac7e547c57b5061240\\memories.mp4

func TrimFileFormat(filePath string) (string, error) {

	formats := []Format{AVI, MP4, MOV}

	for _, format := range formats {
		if strings.Contains(filePath, string(format)) {
			filePath = strings.TrimSuffix(filePath, string(format))
			return filePath, nil
		}
	}

	logrus.Errorf("err_info: invalid file format in %s", filePath)
	return "", fmt.Errorf("err_info: invalid file format in %s", filePath)

}
