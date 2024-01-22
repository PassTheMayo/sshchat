package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"golang.org/x/crypto/ssh"
)

func LoadPrivateKey(file string) (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile("id_rsa")

	if err != nil {
		return nil, err
	}

	/* publicKeyData, err := os.ReadFile("id_rsa.pub")

	if err != nil {
		return err
	} */

	return UnmarshalPrivateKey(privateKeyData)
}

func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return nil, err
	}

	if err := key.Validate(); err != nil {
		return nil, err
	}

	privateKeyData := MarshalPrivateKeyBytes(key)

	if err := os.WriteFile("id_rsa", privateKeyData, 0777); err != nil {
		return nil, err
	}

	publicKeyData, err := MarshalPublicKeyBytes(&key.PublicKey)

	if err != nil {
		return nil, err
	}

	if err := os.WriteFile("id_rsa.pub", publicKeyData, 0777); err != nil {
		return nil, err
	}

	return key, nil
}

func MarshalPrivateKeyBytes(privateKey *rsa.PrivateKey) []byte {
	data := x509.MarshalPKCS1PrivateKey(privateKey)

	block := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   data,
	}

	return pem.EncodeToMemory(&block)
}

func MarshalPublicKeyBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeySSH, err := ssh.NewPublicKey(publicKey)

	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(publicKeySSH), nil
}

func UnmarshalPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return key, nil
}
