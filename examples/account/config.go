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
	"github.com/joyent/triton-go/network"
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

	nc, err := network.NewClient(config)
	if err != nil {
		log.Fatalf("network.NewClient: %s", err)
	}

	ac, err := account.NewClient(config)
	if err != nil {
		log.Fatalf("account.NewClient: %s", err)
	}

	cfg, err := ac.Config().Get(context.Background(), &account.GetConfigInput{})
	if err != nil {
		log.Fatalf("account.Config.Get: %v", err)
	}
	currentNet := cfg.DefaultNetwork
	fmt.Println("Current Network:", currentNet)

	var defaultNet string
	networks, err := nc.List(context.Background(), &network.ListInput{})
	if err != nil {
		log.Fatalf("network.List: %s", err)
	}
	for _, iterNet := range networks {
		if iterNet.Id != currentNet {
			defaultNet = iterNet.Id
		}
	}
	fmt.Println("Chosen Network:", defaultNet)

	input := &account.UpdateConfigInput{
		DefaultNetwork: defaultNet,
	}
	_, err = ac.Config().Update(context.Background(), input)
	if err != nil {
		log.Fatalf("account.Config.Update: %v", err)
	}

	cfg, err = ac.Config().Get(context.Background(), &account.GetConfigInput{})
	if err != nil {
		log.Fatalf("account.Config.Get: %v", err)
	}
	fmt.Println("Default Network:", cfg.DefaultNetwork)
}
