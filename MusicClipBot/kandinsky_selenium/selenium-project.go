package kandinsky_selenium

import (
	"fmt"
	"os"
	"sai_project/trygpt"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const (
	usersDir = "./users"
)

func (vp *VideoParametrs) CreateClip(driver selenium.WebDriver) error {

	err := vp.GeneratePromptsForKandinsky(driver)

	if err != nil {
		return err
	}

	err = FusionBrainSignUp(driver, vp.FusionBrainEmail, vp.FusionBrainPassword)

	if err != nil {
		return err
	}

	err = driver.Get(fusionBrainAuthURL)

	if err != nil {
		fmt.Printf("Error with GET request to %s for auth: %v\n", fusionBrainAuthURL, err)
		return err
	}

	pagesHandle, err := LoadKandinskyPages(driver, vp.VideoQnty)

	if err != nil {
		return err
	}

	err = vp.GenerateVideos(driver, pagesHandle)

	if err != nil {
		return err
	}

	return nil
}

func (vp *VideoParametrs) GeneratePromptsForKandinsky(driver selenium.WebDriver) error {

	prompt := fmt.Sprintf("Сгенерируй подобный JSON: {scene:[['4 строки'], ['4 строки'], ['4 строки'], ... ]}  Длина массива scene = %d. Генерацию проводи на основе текста песни: %s", vp.VideoQnty, vp.Lyrics)

	for i := 0; i < 7; i++ {

		jsonPrompts, err := trygpt.TryGPTRequset(driver, vp.VideoQnty, vp.CurrerntGPTRespID, prompt)

		if err != nil {
			if strings.Contains(err.Error(), "Response JSON array length not equals") {
				vp.CurrerntGPTRespID += 2
				prompt = fmt.Sprintf("ты сгенрировал не %d запросов", vp.VideoQnty)
				continue
			}
			return err
		}

		vp.Prompts = jsonPrompts["scene"]

		break
	}

	return nil
}

func (vp *VideoParametrs) GenerateVideos(driver selenium.WebDriver, pagesHandle []string) error {

	for i := 0; i < vp.VideoQnty; i++ {

		err := driver.SwitchWindow(pagesHandle[i])

		if err != nil {
			return err
		}

		time.Sleep(2 * time.Second)

		err = GenerateVideo(driver, vp.Prompts[i], vp.IsAnimated)

		if err != nil {
			fmt.Printf("Error on %d page\n", i+1)
			return err
		}

	}

	err := vp.BrowserPagesHandle(driver, pagesHandle)

	if err != nil {
		return err
	}

	err = vp.GetVideos(driver, pagesHandle)

	if err != nil {
		return err
	}

	err = vp.GetTogether()

	if err != nil {
		return err
	}

	return nil
}

func LoadKandinskyPages(driver selenium.WebDriver, videoQnty int) ([]string, error) {

	for i := 0; i < videoQnty-1; i++ {
		_, err := driver.ExecuteScript("window.open('https://fusionbrain.ai/t2v/');", nil)
		if err != nil {
			fmt.Print("Error with load window: ")
			return nil, err
		}

	}

	handle, err := driver.WindowHandles()

	if err != nil {
		fmt.Print("Error with getting window handles: ")
		return nil, err
	}

	fmt.Printf("Handles: %v\n\n", handle)

	return handle, nil
}

func (vp *VideoParametrs) GetVideos(driver selenium.WebDriver, handle []string) error {

	for i := 0; i < vp.VideoQnty; i++ {
		err := driver.SwitchWindow(handle[i])

		if err != nil {
			return err
		}

		kandinskyVideoFilePath := fmt.Sprintf("%s/%s/kandinsky-video %d.mp4", vp.UserDirPath, clipsDir, i)

		kandinskyVideoFile, err := os.Create(kandinskyVideoFilePath)

		if err != nil {
			fmt.Printf("Error creating kandinsky video file by path: %s\n", kandinskyVideoFilePath)
			fmt.Println(err)
			return err
		}

		err = DownloadVideo(driver, kandinskyVideoFile)

		if err != nil {
			return err
		}
	}

	return nil
}
