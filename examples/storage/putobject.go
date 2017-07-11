package main

import (
	"log"
	"os"

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

	config := &triton.ClientOptions{
		MantaURL:    "https://us-east.manta.joyent.com/",
		AccountName: "tritongo",
		Signers:     []authentication.Signer{sshKeySigner},
	}
	client, err := storage.NewClient(config)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	reader, err := os.Open("foo.txt")
	if err != nil {
		log.Fatalf("os.Open: %s", err)
	}
	defer reader.Close()

	err = client.Objects().Put(&storage.PutObjectInput{
		ObjectPath:   "foo.txt",
		ObjectReader: reader,
	})
	if err != nil {
		log.Fatalf("GetObject(): %s", err)
	}
}
