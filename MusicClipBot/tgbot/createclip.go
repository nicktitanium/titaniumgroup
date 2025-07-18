package tgbot

import (
	"fmt"
	"sai_project/chromedriver"
	"sai_project/kandinsky_selenium"
	music "sai_project/suno_ai"
)

const (
	currentGPTElemIndex = 3
)

func (b *Bot) SongGeneration(chatID, userID int64, prompt string, isInstrumental bool) error {

	sunoData := &music.SongData{}

	songFiles, err := sunoData.GenerateSong(userID, prompt, isInstrumental)

	if err != nil {
		fmt.Print("Error with create song in suno: ")
		return err
	}

	fmt.Println(sunoData)

	b.Users[userID].SongDuration = append(b.Users[userID].SongDuration, sunoData.SongData[0].Duration)
	b.Users[userID].SongDuration = append(b.Users[userID].SongDuration, sunoData.SongData[1].Duration)
	b.Users[userID].Lyrics = sunoData.SongData[0].Lyrics

	fmt.Printf("Music Lyrics:\n\n%s\n\n", sunoData.SongData[0].Lyrics)

	b.Users[userID].State = awaitingSongNumber

	for _, songFile := range songFiles {
		b.sendAudio(songFile, chatID)
	}

	return nil
}

func (b *Bot) CreateVideoClips(userID int64) error {

	driver, err := chromedriver.InitChromeDriver()

	if err != nil {
		return err
	}

	defer driver.Close()

	kandinskyConfData, err := GetConfigData("kandinsky")

	if err != nil {
		return err
	}

	//var videoQnty int
	videoQnty := 7
	// if math.Mod(b.Users[userID].SongDuration[b.Users[userID].SongNumber], 16.0) != 0.0 {
	// 	videoQnty = int(b.Users[userID].SongDuration[b.Users[userID].SongNumber]/16 + 1)
	// } else {
	// 	videoQnty = int(b.Users[userID].SongDuration[b.Users[userID].SongNumber] / 16)
	// }

	clipParametr := kandinsky_selenium.VideoParametrs{
		CurrerntGPTRespID:   currentGPTElemIndex,
		Lyrics:              b.Users[userID].Lyrics,
		VideoQnty:           videoQnty,
		IsAnimated:          b.Users[userID].isAnimated,
		SongNumber:          b.Users[userID].SongNumber,
		UserDirPath:         b.Users[userID].UserDirPath,
		FusionBrainEmail:    kandinskyConfData["email"].(string),
		FusionBrainPassword: kandinskyConfData["password"].(string),
	}

	err = clipParametr.CreateClip(driver)

	if err != nil {
		return err
	}

	return nil
}
