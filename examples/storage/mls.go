package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/storage"
)

func main() {
	keyID := os.Getenv("SDC_KEY_ID")
	accountName := os.Getenv("SDC_ACCOUNT")
	keyPath := os.Getenv("SDC_KEY_FILE")

	privateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Couldn't find key file matching %s\n%s", keyID, err)
	}

	sshKeySigner, err := authentication.NewPrivateKeySigner(keyID, privateKey, accountName)
	if err != nil {
		log.Fatal(err)
	}

	config := &triton.ClientConfig{
		MantaURL:    os.Getenv("MANTA_URL"),
		AccountName: accountName,
		Signers:     []authentication.Signer{sshKeySigner},
	}
	client, err := storage.NewClient(config)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	ctx := context.Background()
	input := &storage.ListDirectoryInput{}
	output, err := client.Dir().List(ctx, input)
	if err != nil {
		log.Fatalf("storage.Dir.List: %s", err)
	}

	spew.Dump(output)
}
