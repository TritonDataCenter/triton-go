package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/storage"
)

func main() {
	keyID := os.Getenv("MANTA_KEY_ID")
	accountName := "tritongo"
	mantaURL := os.Getenv("MANTA_URL")

	sshKeySigner, err := authentication.NewSSHAgentSigner(keyID, os.Getenv("MANTA_USER"))
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

	output, err := client.Objects().Get(context.Background(), &storage.GetObjectInput{
		ObjectPath: "/stor/foo.txt",
	})
	if err != nil {
		log.Fatalf("GetObject(): %s", err)
	}

	defer output.ObjectReader.Close()
	body, err := ioutil.ReadAll(output.ObjectReader)
	if err != nil {
		log.Fatalf("Reading Object: %s", err)
	}

	fmt.Printf("Content-Length: %d\n", output.ContentLength)
	fmt.Printf("Content-MD5: %s\n", output.ContentMD5)
	fmt.Printf("Content-Type: %s\n", output.ContentType)
	fmt.Printf("ETag: %s\n", output.ETag)
	fmt.Printf("Date-Modified: %s\n", output.LastModified.String())
	fmt.Printf("Object:\n\n%s", string(body))
}
