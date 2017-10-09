package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"encoding/pem"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
)

func printAccount(acct *account.Account) {
	fmt.Println("Account ID:", acct.ID)
	fmt.Println("Account Email:", acct.Email)
	fmt.Println("Account Login:", acct.Login)
}

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
		log.Fatalf("compute.NewClient: %s", err)
	}

	acct, err := a.Get(context.Background(), &account.GetInput{})
	if err != nil {
		log.Fatalf("account.Get: %v", err)
	}

	fmt.Println("New ----")
	printAccount(acct)

	input := &account.UpdateInput{
		CompanyName: fmt.Sprintf("%s-old", acct.CompanyName),
	}

	updatedAcct, err := a.Update(context.Background(), input)
	if err != nil {
		log.Fatalf("account.Update: %v", err)
	}

	fmt.Println("New ----")
	printAccount(updatedAcct)
}
