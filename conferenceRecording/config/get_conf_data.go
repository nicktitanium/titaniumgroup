package config

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
)

const configFilePath = "./config.yaml"

func GetConfigData(key string) map[interface{}]interface{} {

	configFileData, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		logrus.Fatalf("err_info: can't reading config file. err_text: %v", err)
	}

	var keyConfig map[interface{}]interface{}

	err = yaml.Unmarshal(configFileData, &keyConfig)

	if err != nil {
		logrus.Fatalf("err_info: can't marshaling config file data. err_text: %v", err)
	}

	return keyConfig[key].(map[interface{}]interface{})

}
