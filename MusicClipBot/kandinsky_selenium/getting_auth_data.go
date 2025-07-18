package kandinsky_selenium

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func (conf *KandinskyConfig) GetAuthData() error {

	file, err := os.Open("config.yaml")

	if err != nil {
		fmt.Print("Error with open telegram config file: ")
	}

	data, err := ioutil.ReadFile(file.Name())

	if err != nil {
		fmt.Println("Error with read telegram config file: ")
		return err
	}

	err = yaml.Unmarshal(data, conf)

	if err != nil {
		fmt.Println("Error with unmarshal data: ")
		return err
	}

	return nil
}
