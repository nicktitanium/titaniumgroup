package kandinsky_selenium

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const fusionBrainAuthURL = "https://fusionbrain.ai/t2v/"

func EmailElem(email string, driver selenium.WebDriver) error {

	emailField, err := driver.FindElement(selenium.ByXPATH, emailFieldXPath)

	if err != nil {

		fmt.Print("Error finding email field: ", err)
		return err
	}

	err = emailField.SendKeys(email)

	if err != nil {
		fmt.Print("Error sending email to field: ")
		return err
	}

	return nil

}

func PasswordFieldElem(password string, driver selenium.WebDriver) error {

	passField, err := driver.FindElement(selenium.ByXPATH, passwordFieldXPath)

	if err != nil {
		fmt.Print("Error finding password field: ")
		return err
	}

	err = passField.SendKeys(password)

	if err != nil {
		fmt.Print("Error sending password to password field: ")
		return err
	}

	return nil
}

func SendAuthData(driver selenium.WebDriver) error {

	signinButton, err := driver.FindElement(selenium.ByXPATH, doAuthButtonXPath)

	if err != nil {
		fmt.Print("Error finding sign in button: ")
		return err
	}

	// Нажимаем на найденную кнопку

	err = signinButton.Click()

	if err != nil {
		fmt.Print("Error clicking on sign in button: ")
		return err
	}

	return nil
}

const (
	makeNewSceneAtemptQnty = 7
)

func GenerateVideo(driver selenium.WebDriver, prompts []string, animated bool) error {

	err := NewScenes(driver)

	if err != nil {
		return err
	}

	//time.Sleep(2 * time.Second)

	err = IsAnimated(driver, animated)

	if err != nil {
		return err
	}

	// err = VideoResolution(driver)

	// if err != nil {
	// 	return err
	// }

	err = FillingScenes(driver, prompts)

	if err != nil {
		return err
	}

	err = StartGeneration(driver)

	if err != nil {
		return err
	}

	return nil
}

// Добавление сцен для генерации видео

func NewScenes(driver selenium.WebDriver) error {

	// Цикл необходим для выявления не загруженной страницы

	var lastError error

	for i := 0; i < makeNewSceneAtemptQnty; i++ {

		// Находим элемент, отвечающий за добавление сцены

		newScene, err := driver.FindElement(selenium.ByXPATH, newSceneButtonXPath)

		if err != nil {

			lastError = err
			fmt.Print("Error with find new scene button: ")

			if _, emailFieldErr := driver.FindElement(selenium.ByXPATH, emailFieldXPath); emailFieldErr != nil {
				return err
			} else if strings.Contains(err.Error(), "no such element") {
				return lastError
			}

			refreshErr := driver.Refresh()

			if refreshErr != nil {
				fmt.Println("Error refresh page creation new scene: ", refreshErr)
				return refreshErr
			}

			time.Sleep(awaitLoadingPage * time.Second)

			continue

		}

		if lastError != nil {
			return lastError
		}

		// Кликаем на кнопку добавления новой сцены дважды, так получится достичь максимально возможного количества сцен

		for j := 0; j < 2; j++ {

			err = newScene.Click()

			if err != nil {
				fmt.Print("Error clicking on new scene button: ")
				return err
			}

			// Подождем, чтобы избежать ошибок

			time.Sleep(2 * time.Second)
		}

		lastError = nil

		break
	}

	return nil
}

func VideoResolution(driver selenium.WebDriver) error {

	video_resolution_botton, err := driver.FindElement(selenium.ByCSSSelector, selectVideoResolutionButtonXPath)
	if err != nil {
		fmt.Print("Error finding selection resolutions button: ", err)
		return err
	}

	err = video_resolution_botton.Click()
	if err != nil {
		fmt.Print("Error clicking on selection resolutions button: ", err)
		return err
	}

	// Выбор разрешения (по умолчанию указан селектор элемента 16:9)

	video_resolution, err := driver.FindElement(selenium.ByCSSSelector, resolutionButtonXPath)
	if err != nil {
		fmt.Print("Error finding resolution button: ", err)
		return err
	}

	err = video_resolution.Click()
	if err != nil {
		fmt.Print("Error clicking on resolution button: ", err)
		return err
	}

	return nil
}

func IsAnimated(driver selenium.WebDriver, animated bool) error {
	if !animated {
		selection, err := driver.FindElement(selenium.ByXPATH, selectVideoModelButtonXPath)

		if err != nil {
			fmt.Print("Error finding selection button: ", err)
			return err
		}

		err = selection.Click()

		if err != nil {
			fmt.Print("Error clicking on selection button: ", err)
			return err
		}

		videoButton, err := driver.FindElement(selenium.ByXPATH, videoModelButtonXpath)

		if err != nil {
			fmt.Print("Error finding video model button: ", err)
			return err
		}

		err = videoButton.Click()

		if err != nil {
			fmt.Print("Error clicking on video model button: ", err)
			return err
		}
	}

	return nil
}

func FillingScenes(driver selenium.WebDriver, prompts []string) error {

	// Пробегаемся по всем 4-ем сценам

	for i := 1; i <= 4; i++ {

		// Генерируем XPath для сцены
		sceneXPath := "/html/body/div/div/div[1]/div[2]/div[2]/div[2]/div/div[2]/div/div/div/div/div[1]/div[" + strconv.Itoa(i) + "]/div/div/div/textarea"

		scence, err := driver.FindElement(selenium.ByXPATH, sceneXPath)
		if err != nil {
			fmt.Print("Error finding prompt's field: ", err)
			return err
		}

		err = scence.Clear()
		if err != nil {
			fmt.Print("Error clearing field: ", err)
			return err
		}

		err = scence.SendKeys(prompts[i-1])
		if err != nil {
			fmt.Printf("Error sending %d prompt to field: %v", i, err)
			return err
		}

	}

	return nil
}

func StartGeneration(driver selenium.WebDriver) error {

	generate, err := driver.FindElement(selenium.ByXPATH, startGenerateButtonXPath)
	if err != nil {
		fmt.Print("Error finding generate button: ", err)
		return err
	}

	err = generate.Click()
	if err != nil {
		fmt.Print("Error clicking on generate button: ", err)
		return err
	}

	time.Sleep(2 * time.Second)

	return nil
}
