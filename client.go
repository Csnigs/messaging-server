package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader *websocket.Upgrader

type client struct {
	// The client id.
	id int

	// The websocket connection.
	ws *websocket.Conn

	// A buffered channel for outbound messages.
	out chan []byte

	// The server struct.
	srv *Server
}

// Read incoming websocket message and forward it to the server message channel.
func (cli *client) reader() {
	for {
		_, rawMsg, err := cli.ws.ReadMessage()
		if err != nil {
			break
		}

		// This whole block should of course not be here and correct message formatting
		// should be handled by the frontend.
		// The reader should do nothing else than passing along the message to the server channel.
		var msg message
		err = json.Unmarshal(rawMsg, &msg)
		if err != nil {
			msg = message{Dest: cli.id, Data: "Wrong message format."}
		} else if msg.Dest == 0 {
			msg = message{Dest: cli.id, Data: "No destination specified."}
		} else if cli.srv.clients[msg.Dest] == nil {
			msg = message{Dest: cli.id, Data: "Unknown destination user."}
		} else {
			msg = message{Src: cli.id, Dest: cli.srv.clients[msg.Dest].id, Data: msg.Data}
		}

		cli.srv.in <- msg
	}
	cli.ws.Close()
}

// Write message from client outbound message channel to the client websocket.
func (cli *client) writer() {
	for msg := range cli.out {
		err := cli.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	cli.ws.Close()
}

type wsHandler struct {
	srv *Server
}

// Spam 2 goroutines for each incoming request: a reader and a writer.
func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := rand.Intn(1000)
	if wsh.srv.clients[id] != nil { // enough for now.
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	cli := &client{id: id, out: make(chan []byte, 256), ws: ws, srv: wsh.srv}
	cli.srv.register <- cli
	defer func() { cli.srv.unregister <- cli }()
	go cli.writer()
	cli.reader()
}
