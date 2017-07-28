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
)

func main() {
	keyID := os.Getenv("SDC_KEY_ID")
	accountName := os.Getenv("SDC_ACCOUNT")
	signer, err := authentication.NewSSHAgentSigner(keyID, accountName)
	if err != nil {
		log.Fatal(err)
	}

	config := &triton.ClientConfig{
		TritonURL:   os.Getenv("SDC_URL"),
		AccountName: accountName,
		Signers:     []authentication.Signer{signer},
	}
	c, err := compute.NewClient(config)
	if err != nil {
		log.Fatalf("NewClient(): %s", err)
	}

	// Create a new instance using our input attributes...
	// https://github.com/joyent/triton-go/blob/master/compute/instances.go#L206
	createInput := &compute.CreateInstanceInput{
		Name:     "go-test1",
		Package:  "g4-highcpu-512M",
		Image:    "1f32508c-e6e9-11e6-bc05-8fea9e979940",
		Networks: []string{"8450dfd7-a150-4c65-b44a-32e06f78ca4d"},
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
		fmt.Println("\nDuration:", time.Since(startTime).String())
		fmt.Println("Name:", instance.Name)
		fmt.Println("State:", instance.State)
	case <-time.After(5 * time.Minute):
		fmt.Println("Timed out")
	}
}
