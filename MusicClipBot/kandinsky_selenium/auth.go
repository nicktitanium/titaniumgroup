package kandinsky_selenium

import (
	"fmt"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

// Авторизация в FusionBrain

const (
	awaitLoadingPage = 5
	authAtemptQnty   = 7
)

func FusionBrainSignUp(driver selenium.WebDriver, email, password string) error {

	var lastError error

	err := driver.Get(fusionBrainAuthURL)

	if err != nil {
		fmt.Printf("Error doing GET request to %s: %v\n", fusionBrainAuthURL, err)
		return err
	}

	time.Sleep(5 * time.Second)

	for i := 0; i < authAtemptQnty; i++ {

		err := EmailElem(email, driver)

		if err != nil {

			fmt.Printf("Error sending email to %s: %v\n", fusionBrainAuthURL, err)
			lastError = err

			if strings.Contains(err.Error(), "no such element") {

				getReqErr := driver.Get(fusionBrainAuthURL)

				if getReqErr != nil {
					fmt.Printf("Error doing GET request to %s: %v\n", fusionBrainAuthURL, getReqErr)
					return lastError
				}

				time.Sleep(awaitLoadingPage * time.Second)

				continue
			}

			return lastError
		}

		lastError = nil

		break

	}

	if lastError != nil {
		fmt.Printf("Can't to do auth on %s\n\nError: %v", fusionBrainAuthURL, lastError)
		return lastError
	}

	// Ищем поле для ввода пароля от аккаунта

	err = PasswordFieldElem(password, driver)

	if err != nil {
		return err
	}

	// Ищем кнопку для входа в аккаунт

	err = SendAuthData(driver)

	if err != nil {
		return err
	}

	time.Sleep(awaitLoadingPage * time.Second)

	return nil

}
