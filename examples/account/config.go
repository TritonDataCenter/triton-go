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
	"github.com/joyent/triton-go/network"
)

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
