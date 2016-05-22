package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMessagingServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MessagingServer Suite")
}
