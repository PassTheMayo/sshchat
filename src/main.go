package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gliderlabs/ssh"

	ssh2 "golang.org/x/crypto/ssh"
)

var (
	key    *rsa.PrivateKey = nil
	config *Config         = &Config{}
	server *ssh.Server     = nil
)

func init() {
	var (
		err error
		t   time.Time = time.Now()
	)

	// Read the configuration data from file
	{
		if err = config.ReadFile("config.yml"); err != nil {
			panic(err)
		}

		log.Printf("Loaded configuration file (%s)\n", time.Since(t).Round(time.Microsecond))
	}

	// Initialize the SSH server
	{
		server = &ssh.Server{
			Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
			Handler: HandleConnection,
			PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
				return true
			},
		}

		log.Println("Initialized SSH server")
	}

	// Load or generate a new RSA key-pair
	{
		t = time.Now()

		if key, err = LoadPrivateKey("id_rsa"); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				panic(err)
			}

			log.Println("RSA key-pair does not exist, generating one for you... (this may take a moment)")

			if key, err = GeneratePrivateKey(); err != nil {
				panic(err)
			}

			log.Printf("Successfully generated RSA key-pair (%s)\n", time.Since(t).Round(time.Millisecond))
		} else {
			log.Printf("Successfully loaded RSA key-pair from file (%s)\n", time.Since(t).Round(time.Millisecond))
		}
	}

	// Convert the RSA key-pair to an SSH signer for use on the server
	{
		signer, err := ssh2.NewSignerFromKey(key)

		if err != nil {
			panic(err)
		}

		server.AddHostKey(signer)
	}
}

func main() {
	log.Printf("Listening on %s:%d\n", config.Host, config.Port)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
