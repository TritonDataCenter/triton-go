package main

import (
	"log"

	"github.com/jen20/manta-go"
	"github.com/jen20/manta-go/authentication"
	"net/http"
	"time"
)

const accountName = "tritongo"

func main() {
	sshKeySigner, err := authentication.NewSSHAgentSigner(
		"fd:9e:9a:9c:28:99:57:05:18:9f:b6:44:6b:cc:fd:3a", accountName)
	if err != nil {
		log.Fatalf("NewSSHAgentSigner: %s", err)
	}

	client, err := manta.NewClient(&manta.ClientOptions{
		Endpoint:    "https://us-east.manta.joyent.com/",
		AccountName: accountName,
		Signers:     []authentication.Signer{sshKeySigner},
	})
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	signed, err := client.SignURL(&manta.SignURLInput{
		ObjectPath:     "books/treasure_island.txt",
		Method:         http.MethodGet,
		ValidityPeriod: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("SignURL: %s", err)
	}

	log.Printf("Signed URL: %s", signed.SignedURL("http"))
}
