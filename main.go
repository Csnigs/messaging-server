package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var (
	configFile = flag.String("config", "", "if specified, uses config file from that path")
	logFile    = flag.String("logfile", "", "if specified, writes log to file")
	verbose    = flag.Bool("verbose", false, "if set, increase the amount of log")

	// Our configuration.
	config Config

	// The homepage html output.
	homeTempl *template.Template
)

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func main() {
	flag.Parse()

	// setup config
	err := SetupConfig(*configFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	// init server
	srv := newserver()
	go srv.run()

	homeTempl = template.Must(template.ParseFiles("client.html"))
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", wsHandler{srv: srv})

	log.Println("Listening on port", config.ListeningPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ListeningPort), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	log.Println("exiting.")
}
