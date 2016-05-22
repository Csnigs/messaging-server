package main

import (
	"fmt"
	"log"
)

type Server struct {
	// The server configuration.
	config *Config

	// Registered client, indexed by their ID.
	clients map[int]*client

	// Inbound messages from the clients connections.
	in chan message

	// Register requests from clients.
	register chan *client

	// Unregister requests from clients.
	unregister chan *client
}

type message struct {
	// Source User ID.
	Src int `json:"src"`

	// Destination User ID.
	Dest int `json:"dest"`

	// Content of the message.
	Data string `json:"data"`
}

func (s *Server) run() {
	for {
		select {
		case cli := <-s.register:
			// Only register client not yet registered.
			if _, present := s.clients[cli.id]; !present {
				if *verbose {
					log.Println("Registering Client:", cli.id)
				}
				s.clients[cli.id] = cli
				cli.out <- []byte(fmt.Sprintf("Connection established, your ID is %d", cli.id))
			} else {
				log.Printf("WARNING: Registration attempt from already registered client %d", cli.id)
			}
		case cli := <-s.unregister:
			if _, present := s.clients[cli.id]; present {
				// Unregister the client.
				if *verbose {
					log.Println("Unregistering Client:", cli.id)
				}
				delete(s.clients, cli.id)
				// Close its outbound channel.
				close(cli.out)
			}
		case msg := <-s.in:
			if *verbose {
				log.Println("Incoming message:", msg)
			}

			cli := s.clients[msg.Dest]
			// safety check: destination user could drop between the frontend check
			// done when the source user sent his message and now.
			if cli == nil {
				continue
			}
			select {
			case cli.out <- []byte(msg.Data):
			default: // If the client channel buffer is full we assume he dropped his connection.
				if *verbose {
					log.Println("Unregistering Client:", cli.id)
				}
				delete(s.clients, cli.id)
				close(cli.out)
			}

		}
	}
}
