package main

import (
	"fmt"
	"log"
)

type server struct {
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

	// Exit channel
	quit chan int
	done bool
}

type message struct {
	// Source User ID.
	Src int `json:"src"`

	// Destination User ID.
	Dest int `json:"dest"`

	// Content of the message.
	Data string `json:"data"`
}

func newserver() *server {
	srv := server{
		clients:    make(map[int]*client),
		in:         make(chan message),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
	return &srv
}

func (s *server) run() {
	for {
		if s.done {
			break
		}

		select {
		case <-s.quit:
			s.done = true
		case cli := <-s.register:
			// Only register client not yet registered.
			if _, present := s.clients[cli.id]; !present {
				log.Println("Registering Client:", cli.id)
				s.clients[cli.id] = cli

				// This out message is a hack and should be handled by the client on successfull connect.
				select {
				case cli.out <- []byte(fmt.Sprintf("Connection established, your ID is %d", cli.id)):
				default:
				}
			} else {
				log.Printf("WARNING: Registration attempt from already registered client %d", cli.id)
			}
		case cli := <-s.unregister:
			if _, present := s.clients[cli.id]; present {
				// Unregister the client.
				log.Println("Unregistering Client:", cli.id)
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
				log.Println("Can't pass message to client, unregistering Client:", cli.id)
				delete(s.clients, cli.id)
				close(cli.out)
			}
		}
	}
	log.Println("Server exiting...")
}
