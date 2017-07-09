package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/compute"
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
		MantaURL:    os.Getenv("MANTA_URL"),
		AccountName: accountName,
		Signers:     []authentication.Signer{sshKeySigner},
	}

	c, err := compute.NewClient(config)
	if err != nil {
		log.Fatalf("compute.NewClient: %s", err)
	}

	listInput := &compute.ListInstancesInput{}
	instances, err := c.Instances().List(context.Background(), listInput)
	if err != nil {
		log.Fatalf("compute.Instances.List: %v", err)
	}
	numInstances := 0
	for _, instance := range instances {
		numInstances++
		fmt.Println(fmt.Sprintf("-- Instance: %v", instance.Name))
	}
}
