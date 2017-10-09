package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"encoding/pem"
	"io/ioutil"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
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
		TritonURL:   os.Getenv("SDC_URL"),
		AccountName: accountName,
		Signers:     []authentication.Signer{signer},
	}

	a, err := account.NewClient(config)
	if err != nil {
		log.Fatalf("failed to init a new account client: %s", err)
	}

	keys, err := a.Keys().List(context.Background(), &account.ListKeysInput{})
	if err != nil {
		log.Fatalf("failed to list keys: %v", err)
	}

	for _, key := range keys {
		fmt.Println("Key Name:", key.Name)
	}

	if key := keys[0]; key != nil {
		input := &account.GetKeyInput{
			KeyName: key.Name,
		}

		key, err := a.Keys().Get(context.Background(), input)
		if err != nil {
			log.Fatalf("failed to get key: %v", err)
		}

		fmt.Println("First Key:", key.Key)
	}
}
