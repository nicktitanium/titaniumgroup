package kandinsky_selenium

import (
	"fmt"
	"time"
	"github.com/tebeka/selenium"
)

// Прослушивание страницы необходимо для:
// 1) Выявления заведомо известных ошибок;
// 2) Поиск сгенерированного видео

const (
	checkingGeneratedVideoInterval  = 20
	awaitingGeneratedVideosInterval = 20
	startGenerateField              = "generate failed"
	invalidPrompt                   = "invalid prompt"
	signOut                         = "sign out"
)

func (vp *VideoParametrs) BrowserPagesHandle(driver selenium.WebDriver, pagesHandle []string) error {

	generatedVideosHandles := make([]string, vp.VideoQnty)
	var generatedVideosQnty int
	timer := time.After(awaitingGeneratedVideosInterval * time.Minute)

	for {
		select {

		case <-timer:
			fmt.Printf("Wait time is over after %d minutes", awaitingGeneratedVideosInterval)
			return fmt.Errorf("Wait time is over after %d minutes", awaitingGeneratedVideosInterval)

		default:

			for i := 0; i < vp.VideoQnty; i++ {

				err := driver.SwitchWindow(pagesHandle[i])

				if err != nil {
					return nil
				}

				_, err = driver.FindElement(selenium.ByXPATH, downloadVideoButtonXPath)

				if err != nil {

					regenerateErr := vp.CrashHandle(driver, WaitCrash(driver), pagesHandle)

					if regenerateErr != nil {
						return regenerateErr
					}
				} else {

					if !IsGeneratedVideoAdded(pagesHandle[i], generatedVideosHandles) {

						generatedVideosHandles = append(generatedVideosHandles, pagesHandle[i])
						generatedVideosQnty++
						fmt.Printf("Video №%d was generated\n", i+1)

						if generatedVideosQnty == vp.VideoQnty {
							return nil
						}
					}
				}
			}

			time.Sleep(checkingGeneratedVideoInterval * time.Second)
		}

	}

}

func WaitCrash(driver selenium.WebDriver) string {

	// Проверка кнопки "Начать генерацию видео", если она есть - ошибка

	_, err := driver.FindElement(selenium.ByCSSSelector, "#app-container > div.styles_root__l_7NR > div.styles_content__WPknI > div.styles_root__OUPZX > div:nth-child(2) > div > div.styles_stat__lt_bk > div.styles_statButton__02yzL > button")
	if err == nil {
		fmt.Println("Trash was found (generate failed): ", err)
		return startGenerateField
	}

	// Проверка не верного промта

	_, err = driver.FindElement(selenium.ByCSSSelector, "#app-container > div.styles_wrap__JLWId.styles_error__Qu1UO.styles_root__9JDuo")
	if err == nil {
		fmt.Println("Trash was found (invalid prompt): ", err)
		return invalidPrompt
	}

	// Проверка выхода из аккаунта

	_, err = driver.FindElement(selenium.ByCSSSelector, "body > div > div.card-pf")
	if err == nil {
		fmt.Println("Trash was found (sign out): ", err)
		return signOut
	}

	return ""
}

func (vp *VideoParametrs) CrashHandle(driver selenium.WebDriver, crash string, pagesHandle []string) error {

	switch crash {

	case startGenerateField:
		err := driver.Refresh()

		if err != nil {
			fmt.Println("Error refreshing page for regenerate: ", err)
			return err
		}

		regenerateErr := vp.GenerateVideos(driver, pagesHandle)

		if regenerateErr != nil {
			return regenerateErr
		}

	case invalidPrompt:
		invalidPromptErr := vp.CreateClip(driver)

		if invalidPromptErr != nil {
			return invalidPromptErr
		}

	case signOut:
		signInErr := FusionBrainSignUp(driver, vp.FusionBrainEmail, vp.FusionBrainPassword)

		if signInErr != nil {
			return signInErr
		}

		regenerateErr := vp.CreateClip(driver)

		if regenerateErr != nil {
			return regenerateErr
		}
	}

	return nil
}

func IsGeneratedVideoAdded(pageHandle string, pagesHandle []string) bool {

	for i := 0; i < len(pagesHandle); i++ {
		if pagesHandle[i] == pageHandle {
			return true
		}
	}

	return false
}
