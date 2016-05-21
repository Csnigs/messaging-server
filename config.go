package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ListeningPort int `json:"listening_port"`
}

const defaultConfigPath = "config/config.json"

func SetupConfig(file string, cfg *Config) error {
	//load config from path
	configFilePath := defaultConfigPath
	if file != "" {
		configFilePath = file
	}

	log.Println("Loading config from", configFilePath)
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(configFile, cfg)
}
