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
	"encoding/pem"
)

func main() {
	keyID := os.Getenv("SDC_KEY_ID")
	accountName := os.Getenv("SDC_ACCOUNT")
	keyMaterial := os.Getenv("SDC_KEY_MATERIAL")

	var signer authentication.Signer
	var err error

	if keyMaterial == "" {
		signer, err = authentication.NewSSHAgentSigner(keyID, accountName)
		if err != nil {
			log.Fatalf("Error Creating SSH Agent Signer: {{err}}", err)
		}
	} else {
		var keyBytes []byte
		if _, err = os.Stat(keyMaterial); err == nil {
			keyBytes, err = ioutil.ReadFile(keyMaterial)
			if err != nil {
				log.Fatalf("Error reading key material from %s: %s",
					keyMaterial, err)
			}
			block, _ := pem.Decode(keyBytes)
			if block == nil {
				log.Fatalf(
					"Failed to read key material '%s': no key found", keyMaterial)
			}

			if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
				log.Fatalf(
					"Failed to read key '%s': password protected keys are\n"+
						"not currently supported. Please decrypt the key prior to use.", keyMaterial)
			}

		} else {
			keyBytes = []byte(keyMaterial)
		}

		signer, err = authentication.NewPrivateKeySigner(keyID, []byte(keyMaterial), accountName)
		if err != nil {
			log.Fatalf("Error Creating SSH Private Key Signer: {{err}}", err)
		}
	}

	config := &triton.ClientConfig{
		MantaURL:    os.Getenv("SDC_URL"),
		AccountName: accountName,
		Signers:     []authentication.Signer{signer},
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
