package main

import (
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader *websocket.Upgrader

type Client struct {
	id   int
	conn *connection
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// A buffered channel for outbound messages.
	out chan []byte

	// The server struct.
	srv *Server
}

// Read incoming websocket message and forward it to the server message channel.
func (c *connection) reader() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.srv.in <- msg
	}
	c.ws.Close()
}

// Write message from client outbound message channel to the client websocket.
func (c *connection) writer() {
	for msg := range c.out {
		err := c.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

type wsHandler struct {
	srv *Server
}

// Spam 2 goroutines for each incoming request: a reader and a writer.
func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{out: make(chan []byte, 256), ws: ws, srv: wsh.srv}

	cli := &Client{id: rand.Intn(1000), conn: c} // Quick ID for now, worst case we won't register the client.
	c.srv.register <- cli
	defer func() { c.srv.unregister <- cli }()
	go c.writer()
	c.reader()
}
