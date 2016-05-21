package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/websocket"
)

type Config struct {
	ListeningPort int       `json:"listening_port"`
	Websocket     WebSocket `json:"ws"`
}

type WebSocket struct {
	ReadBufferSize  int `json:"read_buffer_size"`
	WriteBufferSize int `json:"write_buffer_size"`
}

const defaultConfigPath = "config/config.json"

func SetupConfig(file string, cfg *Config) error {
	//load config from path.
	configFilePath := defaultConfigPath
	if file != "" {
		configFilePath = file
	}

	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configFile, cfg)
	if err != nil {
		return err
	}

	// setup upgrader.
	upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	if cfg.Websocket.ReadBufferSize != 0 {
		upgrader.ReadBufferSize = cfg.Websocket.ReadBufferSize
	}
	if cfg.Websocket.WriteBufferSize != 0 {
		upgrader.WriteBufferSize = cfg.Websocket.WriteBufferSize
	}

	return nil
}
