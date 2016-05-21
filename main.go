package main

import (
	"flag"
	"log"
)

type Server struct {
	configFile string
	logFile    string
	verbose    bool

	Config Config

	clients map[int]*Client
}

func main() {
	// init server
	srv := Server{
		clients: make(map[int]*Client),
	}
	// set flags and config
	flag.StringVar(&srv.configFile, "config", "", "if specified, uses config file from that path")
	flag.StringVar(&srv.logFile, "logfile", "", "if specified, writes log to file")
	flag.Parse()

	cfg := Config{}
	err := SetupConfig(srv.configFile, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// start listening
	log.Println("Signaling server listening on port ", cfg.ListeningPort)

	log.Println("done")
}
