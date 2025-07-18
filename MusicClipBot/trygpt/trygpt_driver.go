package trygpt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tebeka/selenium"
)

const awaitingLoadingPage = 5

func TryGPTRequset(driver selenium.WebDriver, videoQnty, currentElemIndx int, prompt string) (map[string][][]string, error) {

	err := driver.Get("https://trychatgpt.ru/")

	if err != nil {
		fmt.Println("Error GET request to Try GPT: ", err)
		return nil, err
	}

	time.Sleep(awaitingLoadingPage * time.Second)

	err = InputFieldElem(driver, prompt)

	if err != nil {
		return nil, err
	}

	err = SendPromptElem(driver)

	if err != nil {
		return nil, err
	}

	//time.Sleep(15 * time.Second)
	// /html/body/div[1]/div/div[2]/div/div/div/div/div[2]/div[5]/div[2]/pre/code
	// /html/body/div[1]/div/div[2]/div/div/div/div/div[2]/div[9]/div[2]/pre/code
	jsonPrompts, err := ResponseJSON(driver, fmt.Sprintf("/html/body/div[1]/div/div[2]/div/div/div/div/div[2]/div[%d]/div[2]/pre/code", currentElemIndx))

	if err != nil {
		return nil, err
	}

	var PromptsForKandinsky map[string][][]string

	err = json.Unmarshal(jsonPrompts, &PromptsForKandinsky)

	if err != nil {
		fmt.Println("Error unmarshaling Try GPT response data to map: ", err)
		return nil, err
	}

	if len(PromptsForKandinsky["scene"]) != videoQnty {
		fmt.Printf("Response JSON array length not equals %d", videoQnty)
		return nil, fmt.Errorf("Response JSON array length not equals %d", videoQnty)
	}

	//time.Sleep(15 * time.Second)

	return PromptsForKandinsky, err
}

func StartDriver() (selenium.WebDriver, error) {
	cap := selenium.Capabilities{
		//"--headles"
	}

	_, err := selenium.NewChromeDriverService("./chromedriver.exe", 4444)

	if err != nil {
		fmt.Println("Error running web driver: ", err)
		return nil, err
	}

	driver, err := selenium.NewRemote(cap, "")

	if err != nil {
		fmt.Println("Error start new remote: ", err)
		return nil, err
	}

	err = driver.MaximizeWindow("")

	if err != nil {
		fmt.Println("Error maximazing browser window: ", err)
		return nil, err
	}

	return driver, nil
}

// deleteAdverismentButton, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/div[2]")

// if err != nil {
// 	fmt.Println("Error finding delete advertisment button: ", err)
// 	return
// }

// err = deleteAdverismentButton.Click()

// if err != nil {
// 	fmt.Println("Error clicking on delete advertisment button: ", err)
// 	return
// }
