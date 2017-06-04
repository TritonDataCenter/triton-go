package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/jen20/manta-go"
	"github.com/jen20/manta-go/authentication"
)

func main() {
	sshKeySigner, err := authentication.NewSSHAgentSigner(
		"fd:9e:9a:9c:28:99:57:05:18:9f:b6:44:6b:cc:fd:3a", "tritongo")
	if err != nil {
		log.Fatalf("NewSSHAgentSigner: %s", err)
	}

	client, err := manta.NewClient(&manta.ClientOptions{
		Endpoint:    "https://us-east.manta.joyent.com/",
		AccountName: "tritongo",
		Signers:     []authentication.Signer{sshKeySigner},
	})
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	output, err := client.ListDirectory(&manta.ListDirectoryInput{})
	if err != nil {
		log.Fatalf("ListDirectory(): %s", err)
	}

	spew.Dump(output)
}
