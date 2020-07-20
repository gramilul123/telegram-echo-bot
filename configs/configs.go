package configs

import (
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	URL     string `yaml:"url"`
	Token   string `yaml:"token"`
	LogFile string `yaml:"logfile"`
	Port    string `yaml:"port"`
}

var clientConfig *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		yamlFile, err := os.Open("./configs/configs.yml")
		if err != nil {
			panic(err)
		}
		defer yamlFile.Close()

		byteValue, err := ioutil.ReadAll(yamlFile)
		if err != nil {
			panic(err)
		}

		var configData Config
		err = yaml.Unmarshal(byteValue, &configData)
		if err != nil {
			panic(err)
		}

		clientConfig = &configData

		if clientConfig.Port == "" {
			clientConfig.Port = os.Getenv("PORT")
		}
	})

	return clientConfig
}
