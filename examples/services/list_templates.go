package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/services"
)

func main() {
	keyID := os.Getenv("TRITON_KEY_ID")
	accountName := os.Getenv("TRITON_ACCOUNT")
	keyMaterial := os.Getenv("TRITON_KEY_MATERIAL")
	userName := os.Getenv("TRITON_USER")

	var signer authentication.Signer
	var err error

	if keyMaterial == "" {
		input := authentication.SSHAgentSignerInput{
			KeyID:       keyID,
			AccountName: accountName,
			Username:    userName,
		}
		signer, err = authentication.NewSSHAgentSigner(input)
		if err != nil {
			log.Fatalf("Error Creating SSH Agent Signer: %v", err)
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

		input := authentication.PrivateKeySignerInput{
			KeyID:              keyID,
			PrivateKeyMaterial: keyBytes,
			AccountName:        accountName,
			Username:           userName,
		}
		signer, err = authentication.NewPrivateKeySigner(input)
		if err != nil {
			log.Fatalf("Error Creating SSH Private Key Signer: %v", err)
		}
	}

	config := &triton.ClientConfig{
		ServicesURL: "http://localhost:3000/",
		AccountName: accountName,
		Username:    userName,
		Signers:     []authentication.Signer{signer},
	}

	s, err := services.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create new service groups client: %v", err)
	}

	input := &services.ListTemplatesInput{}
	templates, err := s.Templates().List(context.Background(), input)
	if err != nil {
		log.Fatalf("failed to list machine templates: %v", err)
	}
	for _, template := range templates {
		fmt.Println("Template Name", template.Name)
	}
}