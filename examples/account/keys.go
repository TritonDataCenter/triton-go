package main

import (
	"context"
	"fmt"
	"log"
	"os"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
)

func main() {
	keyID, foundKey := os.LookupEnv("SDC_KEY_ID")
	if !foundKey {
		log.Fatal("Couldn't find \"SDC_KEY_ID\" in your environment")
	}

	accountName, foundAccount := os.LookupEnv("SDC_ACCOUNT")
	if !foundAccount {
		log.Fatal("Couldn't find \"SDC_ACCOUNT\" in your environment")
	}

	tritonURL, foundURL := os.LookupEnv("SDC_URL")
	if !foundURL {
		log.Fatal("Couldn't find \"SDC_URL\" in your environment")
	}

	signer, err := authentication.NewSSHAgentSigner(keyID, accountName)
	if err != nil {
		log.Fatal(err)
	}

	config := &triton.ClientConfig{
		TritonURL:   tritonURL,
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
