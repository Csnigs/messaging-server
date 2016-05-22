package main

import (
	"bytes"
	"fmt"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Messaging server", func() {

	var (
		srv *server
		cli *client
		//otherCli *client
		id     = 1
		msg    message
		logs   bytes.Buffer
		cliMsg []byte
	)

	BeforeEach(func() {
		log.SetOutput(&logs)

		srv = newserver()
		go srv.run()

		cli = &client{id: id, out: make(chan []byte, 256), srv: srv}
	})

	AfterEach(func() {
		logs.Truncate(0)
	})

	Describe("on registration", func() {

		BeforeEach(func() {
			srv.register <- cli
		})

		It("should add the client to the registry.", func() {
			registeredCli := srv.clients[id]
			Expect(*registeredCli).To(Equal(*cli))
			Expect(logs.String()).Should(ContainSubstring(fmt.Sprintf("Registering Client: %d", id)))
		})

		It("should send acknowledgment message to client.", func() {
			Expect(string(<-cli.out)).To(Equal(fmt.Sprintf("Connection established, your ID is %d", id)))
		})
	})

	Describe("on unregistration", func() {

		It("should remove the client from the registry.", func() {
			srv.register <- cli
			srv.unregister <- cli
			// Send an empty message to block the test suite and give the server time to unregister the cli.
			srv.in <- message{}
			Expect(logs.String()).Should(ContainSubstring(fmt.Sprintf("Unregistering Client: %d", id)))
			Expect(srv.clients[id]).Should(BeNil())
		})
	})

	Describe("on message", func() {

		BeforeEach(func() {
			srv.register <- cli
			<-cli.out
		})

		It("Should transmit it if the destination user is known", func() {
			msg = message{Dest: id, Data: "hello"}
			srv.in <- msg
			Expect(string(<-cli.out)).To(Equal(fmt.Sprintf("hello")))
		})

		It("Should not transmit it if the destination user is unknown", func() {
			msg = message{Dest: 2, Data: "hello"}
			srv.in <- msg

			select {
			case cliMsg = <-cli.out:
			default:
				cliMsg = []byte("No message")
			}
			Expect(string(cliMsg)).To(Equal(fmt.Sprintf("No message")))
		})

		It("Should unregister the user if channel is not receiving", func() {
			close(cli.out)
			cli.out = make(chan []byte, 0)
			msg = message{Dest: 1, Data: "hello"}
			srv.in <- msg
			srv.in <- message{}
			Expect(logs.String()).Should(ContainSubstring(fmt.Sprintf("Can't pass message to client, unregistering Client: %d", id)))
			Expect(srv.clients[id]).Should(BeNil())
		})
	})
})
