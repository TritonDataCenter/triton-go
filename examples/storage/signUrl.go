package main

import (
	"log"
	"os"

	"net/http"
	"time"

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

	signed, err := client.SignURL(&storage.SignURLInput{
		ObjectPath:     "/stor/foo.txt",
		Method:         http.MethodGet,
		ValidityPeriod: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("SignURL: %s", err)
	}

	log.Printf("Signed URL: %s", signed.SignedURL("http"))
}
