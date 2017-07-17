package main

import (
	"io/ioutil"
	"log"
	"os"

	"net/http"
	"time"

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

	input := &storage.SignURLInput{
		ObjectPath:     "/stor/books/treasure_island.txt",
		Method:         http.MethodGet,
		ValidityPeriod: 5 * time.Minute,
	}
	signed, err := client.SignURL(input)
	if err != nil {
		log.Fatalf("SignURL: %s", err)
	}

	log.Printf("Signed URL: %s", signed.SignedURL("http"))
}
