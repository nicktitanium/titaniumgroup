package suno_ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	aceDataURL   = "https://api.acedata.cloud/suno/audios"
	configFile   = "./config.yaml"
	usersDir     = "./users/"
	userSongsDir = "/songs"
)

func (sd *SongData) GenerateSong(userID int64, prompt string, isInstrumental bool) ([]*os.File, error) {

	aceDataToken, err := GetBotToken()

	if err != nil {
		return nil, err
	}

	err = sd.RequestToAceData(prompt, aceDataToken, isInstrumental)

	if err != nil {
		return nil, err
	}

	songFiles := make([]*os.File, 0, len(sd.SongData))

	for i, song := range sd.SongData {

		songFile, err := DownloadMusic(song.AudioURL, fmt.Sprintf("%s%d%s%s%d.mp3", usersDir, int(userID), userSongsDir, song.Title, i+1))

		//fmt.Println(song.AudioURL)

		if err != nil {
			return nil, err
		}

		songFiles = append(songFiles, songFile)

	}

	// file1, err := os.Open("./songs/Классный с битлз 1.mp3")

	// if err != nil {
	// 	fmt.Println("Error open file 1:")
	// 	return nil, err
	// }

	// defer file1.Close()

	// file2, err := os.Open("./songs/Классный с битлз 2.mp3")

	// if err != nil {
	// 	fmt.Println("Error open file 1:")
	// 	return nil, err
	// }

	// defer file2.Close()

	// songFiles := []*os.File{file1, file2}

	return songFiles, nil

}

func GetBotToken() (string, error) {

	aceDataConf := &AceDataConf{}

	confFile, err := os.Open(configFile)

	if err != nil {
		fmt.Print("Error with open telegram config file: ")
		return "", err
	}

	data, err := ioutil.ReadFile(confFile.Name())

	if err != nil {
		fmt.Println("Error with read telegram config file: ")
		return "", err
	}

	err = yaml.Unmarshal(data, aceDataConf)

	if err != nil {
		fmt.Println("Error with decoding config data: ")
		return "", err
	}

	return aceDataConf.Data.Token, nil
}

func (sd *SongData) RequestToAceData(prompt, aceDataToken string, isInstrumental bool) error {

	payload := map[string]interface{}{
		"action":         "generate",
		"prompt":         prompt,
		"model":          "chirp-v3-0",
		"continue_at":    0,
		"custom":         false,
		"instrumental":   isInstrumental,
		"style_negative": "",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding payload to JSON data:")
		return err
	}

	fmt.Printf("Payload for Ace Data:\n%v\n\n", string(jsonData))

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, aceDataURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request to AceData service:")
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+aceDataToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error doing request to AceData service:")
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading AceData response body:")
		return err
	}

	fmt.Println(string(body))

	//body := `{"success":true,"data":[{"state":"succeeded","id":"072b23cd-315c-4f68-afff-77a4c5380090","title":"Классный с битлз","image_url":"https://cdn2.suno.ai/image_072b23cd-315c-4f68-afff-77a4c5380090.jpeg","lyric":"[Verse]\nС утра проснулся поздно не беда\nВключаю музыку опять игра\nНа кухне танцую будто бы звезда\nПусть весь мир летит туда-сюда\n\n[Chorus]\nЭй ты смотри на меня\nС ритмом круче битлз всегда\nЗаражаю весельем во всем\nС нами днем и даже ночью\n\n[Verse 2]\nМой друг звонок давай к нам сюда\nСегодня вечеринка будь готова\nРаспахнуты окна счастье не жди\nВсе решит танец начнем дела\n\n[Chorus]\nЭй ты смотри на меня\nС ритмом круче битлз всегда\nЗаражаю весельем во всем\nС нами днем и даже ночью\n\n[Bridge]\nИ пусть не все знают что я звезда\nМой свет зажигает сердца всегда\nЯ просто живу танцую без бед\nСлушай ритм и пополним наш свет\n\n[Chorus]\nЭй ты смотри на меня\nС ритмом круче битлз всегда\nЗаражаю весельем во всем\nС нами днем и даже ночью","audio_url":"https://cdn1.suno.ai/072b23cd-315c-4f68-afff-77a4c5380090.mp3","video_url":"https://cdn1.suno.ai/072b23cd-315c-4f68-afff-77a4c5380090.mp4","created_at":"2024-12-10T13:13:55.319Z","model":"chirp-v3-5","prompt":"Классный с битлз","style":"High quality","duration":146},{"state":"succeeded","id":"816a89f8-2e4f-4689-a4b5-b2068a39d0f2","title":"Классный с битлз","image_url":"https://cdn2.suno.ai/image_816a89f8-2e4f-4689-a4b5-b2068a39d0f2.jpeg","lyric":"[Verse]\nС утра проснулся поздно не беда\nВключаю музыку опять игра\nНа кухне танцую будто бы звезда\nПусть весь мир летит туда-сюда\n\n[Chorus]\nЭй ты смотри на меня\nС ритмом круче битлз всегда\nЗаражаю весельем во всем\nС нами днем и даже ночью\n\n[Verse 2]\nМой друг звонок давай к нам сюда\nСегодня вечеринка будь готова\nРаспахнуты окна счастье не жди\nВсе решит танец начнем дела\n\n[Chorus]\nЭй ты смотри на меня\nС ритмом круче битлз всегда\nЗаражаю весельем во всем\nС нами днем и даже ночью\n\n[Bridge]\nИ пусть не все знают что я звезда\nМой свет зажигает сердца всегда\nЯ просто живу танцую без бед\nСлушай ритм и пополним наш свет\n\n[Chorus]\nЭй ты смотри на меня\nС ритмом круче битлз всегда\nЗаражаю весельем во всем\nС нами днем и даже ночью","audio_url":"https://cdn1.suno.ai/816a89f8-2e4f-4689-a4b5-b2068a39d0f2.mp3","video_url":"https://cdn1.suno.ai/816a89f8-2e4f-4689-a4b5-b2068a39d0f2.mp4","created_at":"2024-12-10T13:13:55.264Z","model":"chirp-v3-5","prompt":"Классный с битлз","style":"High quality","duration":199}],"task_id":"63628e1e-c40f-4ee4-bf1a-a5a48f02d22d"}`

	err = json.Unmarshal([]byte(body), &sd)

	if err != nil {
		fmt.Println("Error decoding AceData response body in struct:")
		return err
	}

	return nil
}

func DownloadMusic(url string, filepath string) (*os.File, error) {

	if len(filepath) < 4 || filepath[len(filepath)-4:] != ".mp3" {
		filepath += ".mp3"
	}

	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error request: ", resp.StatusCode)
		return nil, err
	}

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return nil, err
	}

	return file, nil

}
