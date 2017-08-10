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

func printAccount(acct *account.Account) {
	fmt.Println("Account ID:", acct.ID)
	fmt.Println("Account Email:", acct.Email)
	fmt.Println("Account Login:", acct.Login)
}

func main() {
	keyID := os.Getenv("SDC_KEY_ID")
	accountName := os.Getenv("SDC_ACCOUNT")

	sshKeySigner, err := authentication.NewSSHAgentSigner(keyID, accountName)
	if err != nil {
		log.Fatalf("NewSSHAgentSigner: %s", err)
	}

	client, err := account.NewClient(&triton.ClientConfig{
		TritonURL:   os.Getenv("SDC_URL"),
		AccountName: accountName,
		Signers:     []authentication.Signer{sshKeySigner},
	})
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	acct, err := client.Get(context.Background(), &account.GetInput{})
	if err != nil {
		log.Fatalf("account.Get: %v", err)
	}

	printAccount(acct)
}
