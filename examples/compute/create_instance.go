package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/network"
)

const (
	PackageName  = "g4-highcpu-512M"
	ImageName    = "ubuntu-16.04"
	ImageVersion = "20170403"
	NetworkName  = "Joyent-SDC-Private"
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
	c, err := compute.NewClient(config)
	if err != nil {
		log.Fatalf("Compute NewClient(): %s", err)
	}
	n, err := network.NewClient(config)
	if err != nil {
		log.Fatalf("Network NewClient(): %s", err)
	}

	images, err := c.Images().List(context.Background(), &compute.ListImagesInput{
		Name:    ImageName,
		Version: ImageVersion,
	})
	img := images[0]

	var net *network.Network
	nets, err := n.List(context.Background(), &network.ListInput{})
	if err != nil {
		log.Fatalf("Network List(): %s", err)
	}
	for _, found := range nets {
		if found.Name == NetworkName {
			net = found
		}
	}

	// Create a new instance using our input attributes...
	// https://github.com/joyent/triton-go/blob/master/compute/instances.go#L206
	createInput := &compute.CreateInstanceInput{
		Name:     "go-test1",
		Package:  PackageName,
		Image:    img.ID,
		Networks: []string{net.Id},
		Metadata: map[string]string{
			"user-script": "<your script here>",
		},
		Tags: map[string]string{
			"tag1": "value1",
		},
		CNS: compute.InstanceCNS{
			Services: []string{"frontend", "web"},
		},
	}
	startTime := time.Now()
	created, err := c.Instances().Create(context.Background(), createInput)
	if err != nil {
		log.Fatalf("Create(): %v", err)
	}

	// Wait for provisioning to complete...
	state := make(chan *compute.Instance, 1)
	go func(createdID string, c *compute.ComputeClient) {
		for {
			time.Sleep(1 * time.Second)
			instance, err := c.Instances().Get(context.Background(), &compute.GetInstanceInput{
				ID: createdID,
			})
			if err != nil {
				log.Fatalf("Get(): %v", err)
			}
			if instance.State == "running" {
				state <- instance
			} else {
				fmt.Print(".")
			}
		}
	}(created.ID, c)

	select {
	case instance := <-state:
		fmt.Printf("\nDuration: %s\n", time.Since(startTime))
		fmt.Println("Name:", instance.Name)
		fmt.Println("State:", instance.State)
	case <-time.After(5 * time.Minute):
		fmt.Println("Timed out")
	}
}
