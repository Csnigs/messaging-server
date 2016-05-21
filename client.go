package main

import "golang.org/x/net/websocket"

type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
}
