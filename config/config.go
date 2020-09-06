package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Path string `json:"path"`
	Port int    `json:"port"`
}

var ServerConfig Config

func ReadConfig(path string) error {
	configFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer configFile.Close()
	bytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &ServerConfig)
	return err
}
