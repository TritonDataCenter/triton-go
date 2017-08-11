package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/network"
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
	n, err := network.NewClient(config)
	if err != nil {
		log.Fatalf("Network NewClient(): %s", err)
	}

	fabric, err := n.Fabrics().Create(context.Background(), &network.CreateFabricInput{
		FabricVLANID:     2,
		Name:             "testnet",
		Description:      "This is a test network",
		Subnet:           "10.50.1.0/24",
		ProvisionStartIP: "10.50.1.10",
		ProvisionEndIP:   "10.50.1.240",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Fabric was successfully created!")
	fmt.Println("Name:", fabric.Name)
	time.Sleep(5 * time.Second)

	err = n.Fabrics().Delete(context.Background(), &network.DeleteFabricInput{
		FabricVLANID: 2,
		NetworkID:    fabric.Id,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Fabric was successfully deleted!")
	time.Sleep(5 * time.Second)

	fwrule, err := n.Firewall().CreateRule(context.Background(), &network.CreateRuleInput{
		Enabled: false,
		Rule:    "FROM any TO tag \"bone-thug\" = \"basket-ball\" ALLOW udp PORT 8600",
	})

	fmt.Println("Firewall Rule was successfully added!")
	time.Sleep(5 * time.Second)

	err = n.Firewall().DeleteRule(context.Background(), &network.DeleteRuleInput{
		ID: fwrule.ID,
	})

	fmt.Println("Firewall Rule was successfully deleted!")
	time.Sleep(5 * time.Second)

}
