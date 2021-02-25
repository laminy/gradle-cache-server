package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Path          string `json:"path"`
	Port          int    `json:"port"`
	Scan          string `json:"scan"`
	ScanInterval  time.Duration
	Alive         string `json:"alive"`
	AliveInterval time.Duration
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
	ServerConfig.ScanInterval = parseDurationString(ServerConfig.Scan)
	ServerConfig.AliveInterval = parseDurationString(ServerConfig.Alive)
	return err
}

func parseDurationString(value string) (duration time.Duration) {
	valueLen := len(value)
	if valueLen < 2 {
		log.Fatal("invalid value format", value)
	}
	timeMod := value[valueLen-1:]
	switch timeMod {
	case "s":
		duration = time.Second
		break
	case "m":
		duration = time.Minute
		break
	case "h":
		duration = time.Hour
		break
	case "d":
		duration = time.Hour * 24
		break
	case "w":
		duration = time.Hour * 24 * 7
		break
	default:
		log.Fatal("unknown interval format", timeMod)
	}
	valueMod := value[0 : valueLen-1]
	iVal, err := strconv.ParseInt(valueMod, 10, 32)
	if err != nil {
		log.Fatal("invalid number format", err)
	}
	duration *= time.Duration(iVal)
	return duration
}
