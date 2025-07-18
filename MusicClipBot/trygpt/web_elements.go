package trygpt

import (
	"fmt"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const (
	inputFieldXPath       = "/html/body/div[1]/div/div[2]/div/div/div/div/div[3]/div/div/textarea"
	sendPromptButtonXPath = "/html/body/div[1]/div/div[2]/div/div/div/div/div[3]/div/div/button"
)

func InputFieldElem(driver selenium.WebDriver, prompt string) error {

	inputField, err := driver.FindElement(selenium.ByXPATH, inputFieldXPath)

	if err != nil {
		fmt.Println("Error finding input field: ", err)
		return err
	}

	err = inputField.Clear()

	if err != nil {
		fmt.Println("Error clearing input field: ", err)
		return err
	}

	err = inputField.SendKeys(prompt)

	if err != nil {
		fmt.Println("Error sending prompt to input field: ", err)
		return err
	}

	return nil
}

func SendPromptElem(driver selenium.WebDriver) error {

	sendPromptButton, err := driver.FindElement(selenium.ByXPATH, sendPromptButtonXPath)

	if err != nil {
		fmt.Println("Error finding send button: ", err)
		return err
	}

	err = sendPromptButton.Click()

	if err != nil {
		fmt.Println("Error clicking on send button: ", err)
		return err
	}

	return nil
}

func ResponseJSON(driver selenium.WebDriver, currentXPath string) ([]byte, error) {

	awaitingJSON := func(driver selenium.WebDriver) (bool, error) {

		_, err := driver.FindElement(selenium.ByXPATH, currentXPath)

		if err != nil {

			// Если элемент не найден, значит, вероятно, он еще не прогрузился и ответ от нейросети еще не был сгенерирован

			if strings.Contains(strings.ToLower(err.Error()), "no such element") {
				return false, nil
			}

			return false, err
		}

		return true, nil
	}

	err := driver.WaitWithTimeout(awaitingJSON, 30*time.Second)

	if err != nil {
		fmt.Println("Error awaiting JSON data from Try GPT: ")
		return nil, err
	}

	jsonData, err := driver.FindElement(selenium.ByXPATH, currentXPath)

	if err != nil {
		fmt.Println("Error finding JSON response data: ", err)
		return nil, err
	}

	jsonRequestToFusionBrain, err := jsonData.Text()

	if err != nil {
		fmt.Println("Error getting JSON response data text: ", err)
		return nil, err
	}

	return []byte(jsonRequestToFusionBrain), nil
}
