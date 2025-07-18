package kandinsky_selenium

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/tebeka/selenium"
)

// Скачивание сгенерированного видео с вкладки

func DownloadVideo(driver selenium.WebDriver, kandinskyVideoFile *os.File) error {

	elem, err := driver.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div[2]/div[1]/div[2]/div[2]/div/div/div/div/video")

	if err != nil {
		fmt.Print("Error with find download button: ")
		return err
	}

	// Берем атрибут элемента

	src, err := elem.GetAttribute("src")
	if err != nil {
		fmt.Print("Error getting attribute:", err)
		return err
	}

	// Проверка, начинается ли строка с ожидаемого префикса

	if !strings.HasPrefix(src, "data:video/mp4;base64,") {
		return fmt.Errorf("Expected base64 video format, got: %s", src)
	}

	// Извлекаем только Base64 данные, убираем префикс

	base64Data := strings.TrimPrefix(src, "data:video/mp4;base64,")

	// Декодируем Base64 данные

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Print("Error with decoding base64 data: ")
		return err
	}

	// Сохраняем декодированные данные в файл

	err = os.WriteFile(kandinskyVideoFile.Name(), data, 0644)
	if err != nil {
		fmt.Println("Error writing data to file: ", err)
		return err
	}

	fmt.Println("File successfully download: ", kandinskyVideoFile.Name())

	return nil

}
