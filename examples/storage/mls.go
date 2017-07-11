package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/storage"
)

func main() {
	const accountName = os.Getenv("MANTA_USER")

	sshKeySigner, err := authentication.NewSSHAgentSigner(
		"fd:9e:9a:9c:28:99:57:05:18:9f:b6:44:6b:cc:fd:3a", accountName)
	if err != nil {
		log.Fatalf("NewSSHAgentSigner: %s", err)
	}

	config := &triton.ClientConfig{
		MantaURL:    "https://us-east.manta.joyent.com/",
		AccountName: accountName,
		Signers:     []authentication.Signer{sshKeySigner},
	}
	client, err := storage.NewClient(config)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	output, err := client.Dir().List(&storage.ListDirectoryInput{})
	if err != nil {
		log.Fatalf("ListDirectory(): %s", err)
	}

	spew.Dump(output)
}
