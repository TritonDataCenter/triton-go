package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jen20/triton-go"
	"github.com/jen20/triton-go/authentication"
	"log"
)

func main() {
	sshKeySigner, err := authentication.NewSSHAgentSigner(
		"1b:bc:29:48:89:af:72:63:f0:83:b8:11:b6:4d:ff:3f", "hashicorp")
	if err != nil {
		log.Fatalf("NewSSHAgentSigner: %s", err)
	}

	client, err := triton.NewClient("https://us-sw-1.api.joyent.com/", "hashicorp", sshKeySigner)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	keys, err := client.Keys().ListKeys()
	if err != nil {
		log.Fatalf("ListKeys(): %s", err)
	}

	spew.Dump(keys)
}
