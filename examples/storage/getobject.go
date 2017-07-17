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

	obj, err := client.Objects().Get(context.Background(), &storage.GetObjectInput{
		ObjectPath: "/stor/books/dracula.txt",
	})
	if err != nil {
		log.Fatalf("compute.Objects.Get: %s", err)
	}

	body, err := ioutil.ReadAll(obj.ObjectReader)
	if err != nil {
		log.Fatalf("compute.Objects.Get: %s", err)
	}
	defer obj.ObjectReader.Close()

	fmt.Printf("Content-Length: %d\n", obj.ContentLength)
	fmt.Printf("Content-MD5: %s\n", obj.ContentMD5)
	fmt.Printf("Content-Type: %s\n", obj.ContentType)
	fmt.Printf("ETag: %s\n", obj.ETag)
	fmt.Printf("Date-Modified: %s\n", obj.LastModified.String())
	fmt.Printf("Length: %d\n", len(body))
}
