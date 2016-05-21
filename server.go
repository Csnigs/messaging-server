package main

import "log"

type Server struct {
	// The server configuration.
	config *Config

	// Registered client, indexed by their ID.
	clients map[int]*Client

	// Inbound messages from the clients connections.
	in chan []byte

	// Register requests from clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
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
				close(cli.conn.out)
			}
		case msg := <-s.in:
			if *verbose {
				log.Println("Incoming message:", string(msg))
			}

			// Broadcast the message to all client for now.
			for _, cli := range s.clients {
				select {
				case cli.conn.out <- msg:
				default: // If the client channel buffer is full we assume he dropped his connection.
					delete(s.clients, cli.id)
					close(cli.conn.out)
				}
			}
		}
	}
}
