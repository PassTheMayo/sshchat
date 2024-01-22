package main

import (
	"log"

	"github.com/gliderlabs/ssh"
)

func HandleConnection(session ssh.Session) {
	log.Printf("Received a connection from %s\n", session.RemoteAddr())

	session.Write([]byte("Hello, world!\n"))

	defer session.Close()
}
