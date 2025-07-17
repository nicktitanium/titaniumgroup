package webui

import (
	llmstudioapi "conference_recognition/conference-summary/llm-studio-api"
	whisperuiapi "conference_recognition/conference-summary/whisper-ui-api"
	"conference_recognition/config"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

type SelfHostData struct {
	BaseUrl          string
	UserAudioDir     string
	UserSubtitlesDir string
	UserSummaryDir   string
	HTMLPath         string
}

func getHostData() SelfHostData {

	confData := config.GetConfigData("self_host")

	return SelfHostData{
		BaseUrl:          confData["url"].(map[interface{}]interface{})["base_url"].(string),
		UserAudioDir:     confData["path"].(map[interface{}]interface{})["audio_dir"].(string),
		UserSubtitlesDir: confData["path"].(map[interface{}]interface{})["subtitles_dir"].(string),
		UserSummaryDir:   confData["path"].(map[interface{}]interface{})["summary_dir"].(string),
		HTMLPath:         confData["path"].(map[interface{}]interface{})["html"].(string),
	}

}

func homePageHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		fmt.Fprint(w, `
			<html>
				<head>
					<meta charset="UTF-8">
				</head>
			<body>
				<form action="/generate" method="post" enctype="multipart/form-data">
				<label for="video">Выберите видео файл</label>
				<input type="file" id="video" name="video" accept="video/*" required>
				<br/>
				<br/>
				<label for="prompt">Запрос для нейросети</label>
				<input type="text" id="prompt" name="prompt">
				<br/>
				<br/>
				<label for="temperature">Процент отклонения от текста конференции(max = 1)</label>
				<input type="text" id="temperature" name="temperature">		
				<br/>
				<br/>
				<button type="submit">Загрузить и сгенерировать</button>
				</form>
			</body>
			</html>`)

		return

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {

	confData := getHostData()

	switch r.Method {

	case http.MethodPost:

		newAudioFile, err := sendAudioToDir(r, confData.UserAudioDir)

		if err != nil {
			http.Error(w, "Somethig went wrong. Please call the administrator or try to generate again", http.StatusInternalServerError)
			return
		}

		prompt := r.FormValue("prompt")
		temperature := r.FormValue("temperature")

		var temperatureNumber float32

		if intNumber, err := strconv.Atoi(temperature); err == nil && intNumber <= 1 && intNumber > 0 {

			temperatureNumber = float32(intNumber)

		} else if floatNumber, err := strconv.ParseFloat(temperature, 32); err == nil && floatNumber <= 1 && floatNumber > 0 {

			temperatureNumber = float32(floatNumber)

		} else if floatNumber == 0.0 || intNumber == 0 {

			temperatureNumber = 0.0

		} else {
			http.Error(w, "Invalid temperature value", http.StatusBadRequest)
			return

		}

		logrus.Infof("start creating subtitles for %s, payload: %s, temperature: %.1f", filepath.Base(newAudioFile.Name()), prompt, temperatureNumber)

		transcriptFilePath, err := whisperuiapi.CreateSubtitles(newAudioFile.Name(), temperatureNumber)

		if err != nil {
			http.Error(w, "Error Creating Subtitles", http.StatusInternalServerError)
			return
		}

		fileName := filepath.Base(transcriptFilePath)

		sendSubtitlesToDir(fmt.Sprintf("%s/%s", confData.UserSubtitlesDir, fileName), transcriptFilePath)

		logrus.Info("start creating summary")

		conferenceSummary, err := llmstudioapi.CreateSummary(transcriptFilePath, prompt)

		if err != nil {
			http.Error(w, "Somethig went wrong. Please call the administrator or try to generate again", http.StatusInternalServerError)
			return
		}

		sendSummaryToDir(fmt.Sprintf("%s/%s", confData.UserSummaryDir, fileName), conferenceSummary)

		content, err := ioutil.ReadFile(transcriptFilePath)

		if err != nil {
			logrus.Errorf("err_info: can't read transcript file %s, err_text: %v", transcriptFilePath, err)
			http.Error(w, "Somethig went wrong. Please call the administrator or try to generate again", http.StatusInternalServerError)
			return
		}

		//encodedURL := url.QueryEscape(transcriptFilePath)

		response := fmt.Sprintf(`
		   <html>
			<head>
				<meta charset="UTF-8">
			</head>
		   <body>
		   		<h2>Название конфиренции</h2>
				<p>%s</p>
		        <h2>Сгенерированный текст:</h2>
		        <p>%s</p>
				<br/>
				<h2>Субтитры</h2>
				<p>%s</p>
			</body>
		   </html>`, fileName, conferenceSummary, content)

		w.Header().Set("Content-Type", "text/html")

		if _, err = w.Write([]byte(response)); err != nil {
			logrus.Errorf("err_info: can't send html markup (summary html), err_text: %v", err)
			http.Error(w, "Somethig went wrong. Please call the administrator or try to generate again", http.StatusInternalServerError)
			return
		}

		logrus.Info("summary successfully sent!")

		return
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

}

func sendAudioToDir(r *http.Request, userAudioDir string) (*os.File, error) {

	conferenceVideo, conferenceVideoInfo, err := r.FormFile("video")

	if err != nil {
		logrus.Errorf("err_info: can't read from html form, err_text: %v", err)
		return nil, err
	}

	defer conferenceVideo.Close()

	newAudioFile, err := os.Create(fmt.Sprintf("%s/%s", userAudioDir, conferenceVideoInfo.Filename))

	if err != nil {
		logrus.Fatalf("info: can't create new audio file for recognition, file_name: %s, err_text: %v", conferenceVideoInfo.Filename, err)
	}

	defer newAudioFile.Close()

	if _, err := io.Copy(newAudioFile, conferenceVideo); err != nil {
		logrus.Fatalf("info: can't copy audio file for recognition, file_name: %s, err_text: %v", conferenceVideoInfo.Filename, err)
	}

	return newAudioFile, nil
}

func sendSubtitlesToDir(newFilePath, filePath string) {

	newSubtitlesFile, err := os.Create(newFilePath)

	if err != nil {
		logrus.Fatalf("info: can't create new file for subtitles, file_name: %s, err_text: %v", filepath.Base(newFilePath), err)
	}

	defer newSubtitlesFile.Close()

	subtitlesFile, err := os.Open(filePath)

	if err != nil {
		logrus.Fatalf("info: can't open new subtitles file, file_name: %s, err_text: %v", filepath.Base(filePath), err)
	}

	defer subtitlesFile.Close()

	if _, err = io.Copy(newSubtitlesFile, subtitlesFile); err != nil {
		logrus.Fatalf("info: can't copy subtitles to new user's subtitles file, file_name: %s, err_text: %v", filePath, err)
	}
}

func sendSummaryToDir(filePath, summary string) {

	newAudioFile, err := os.Create(filePath)

	if err != nil {
		logrus.Fatalf("info: can't create new file for subtitles, file_name: %s, err_text: %v", filepath.Base(filePath), err)
	}

	defer newAudioFile.Close()

	if _, err := fmt.Fprint(newAudioFile, summary); err != nil {
		logrus.Fatalf("info: can't write summary to new user's summary file, file_name: %s, err_text: %v", filePath, err)
	}
}
