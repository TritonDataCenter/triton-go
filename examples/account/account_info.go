package main

import (
	"context"
	"fmt"
	"io/ioutil"
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
	keyPath := os.Getenv("SDC_KEY_FILE")

	privateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Couldn't find key file matching %s\n%s", keyID, err)
	}

	sshKeySigner, err := authentication.NewPrivateKeySigner(keyID, privateKey, accountName)
	if err != nil {
		log.Fatal(err)
	}

	config := &triton.ClientConfig{
		TritonURL:   os.Getenv("SDC_URL"),
		AccountName: accountName,
		Signers:     []authentication.Signer{sshKeySigner},
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
