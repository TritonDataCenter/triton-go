//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

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
		TritonURL:   os.Getenv("TRITON_URL"),
		AccountName: accountName,
		Username:    userName,
		Signers:     []authentication.Signer{signer},
	}

	a, err := account.NewClient(config)
	if err != nil {
		log.Fatalf("compute.NewClient: %v", err)
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
