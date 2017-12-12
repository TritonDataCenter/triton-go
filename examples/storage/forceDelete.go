package main

import (
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"

	"context"

	"fmt"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/storage"
)

func main() {
	keyID := os.Getenv("SDC_KEY_ID")
	keyMaterial := os.Getenv("SDC_KEY_MATERIAL")
	mantaUser := os.Getenv("MANTA_USER")
	mantaFolder := os.Getenv("MANTA_FOLDER")

	var signer authentication.Signer
	var err error

	if keyMaterial == "" {
		signer, err = authentication.NewSSHAgentSigner(keyID, mantaUser)
		if err != nil {
			log.Fatalf("Error Creating SSH Agent Signer: %s", err.Error())
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

		signer, err = authentication.NewPrivateKeySigner(keyID, []byte(keyMaterial), mantaUser)
		if err != nil {
			log.Fatalf("Error Creating SSH Private Key Signer: %s", err.Error())
		}
	}

	config := &triton.ClientConfig{
		MantaURL:    os.Getenv("MANTA_URL"),
		AccountName: mantaUser,
		Signers:     []authentication.Signer{signer},
	}

	client, err := storage.NewClient(config)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	err = client.Dir().Delete(context.Background(), &storage.DeleteDirectoryInput{
		DirectoryName: mantaFolder,
		ForceDelete:   true,
	})

	if err != nil {
		log.Fatalf("Error deleting nested folder structure: %s", err.Error())
	}
	fmt.Println("Successfully deleted all nested objects")
}
