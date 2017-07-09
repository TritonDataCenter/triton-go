package main

import (
	"fmt"
	"io/ioutil"
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

	config := &triton.ClientConfig{
		MantaURL:    "https://us-east.manta.joyent.com/",
		AccountName: accountName,
		Signers:     []authentication.Signer{sshKeySigner},
	}
	client, err := storage.NewClient(config)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	output, err := client.Objects().Get(&storage.GetObjectInput{
		ObjectPath: "tempstate.tfstate",
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
