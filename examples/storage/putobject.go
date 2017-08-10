package main

import (
	"context"
	"log"
	"os"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/storage"
)

func main() {
	keyID := os.Getenv("MANTA_KEY_ID")
	accountName := os.Getenv("MANTA_USER")
	mantaURL := os.Getenv("MANTA_URL")

	sshKeySigner, err := authentication.NewSSHAgentSigner(keyID, accountName)
	if err != nil {
		log.Fatalf("NewSSHAgentSigner: %s", err)
	}

	client, err := storage.NewClient(&triton.ClientConfig{
		MantaURL:    mantaURL,
		AccountName: accountName,
		Signers:     []authentication.Signer{sshKeySigner},
	})
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	reader, err := os.Open("foo.txt")
	if err != nil {
		log.Fatalf("os.Open: %s", err)
	}
	defer reader.Close()

	err = client.Objects().Put(context.Background(), &storage.PutObjectInput{
		ObjectPath:   "/stor/foo.txt",
		ObjectReader: reader,
	})
	if err != nil {
		log.Fatalf("GetObject(): %s", err)
	}
}
