package main

import (
	"context"
	"fmt"
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

	output, err := client.Dir().List(context.Background(), &storage.ListDirectoryInput{
		DirectoryName: "/stor",
	})
	if err != nil {
		log.Fatalf("ListDirectory(): %s", err)
	}
	for _, entry := range output.Entries {
		fmt.Println(entry.Name)
	}
}
