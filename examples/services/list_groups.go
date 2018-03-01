package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/services"
)

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
		ServicesURL: "http://localhost:3000/",
		AccountName: accountName,
		Username:    userName,
		Signers:     []authentication.Signer{signer},
	}

	svc, err := services.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create new services client: %v", err)
	}

	fmt.Println("---")

	listInput := &services.ListGroupsInput{}
	groups, err := svc.Groups().List(context.Background(), listInput)
	if err != nil {
		log.Fatalf("failed to list service groups: %v", err)
	}

	for _, grp := range groups {
		fmt.Printf("Group ID: %v\n", grp.ID)
		fmt.Printf("Group Name: %v\n", grp.GroupName)
		fmt.Printf("Group TemplateID: %v\n", grp.TemplateID)
		fmt.Printf("Group AccountID: %v\n", grp.AccountID)
		fmt.Printf("Group Capacity: %v\n", grp.Capacity)
		fmt.Printf("Group HealthCheckInterval: %v\n", grp.HealthCheckInterval)
		fmt.Println("")
	}

	fmt.Println("---")

	listTmpls := &services.ListTemplatesInput{}
	templates, err := svc.Templates().List(context.Background(), listTmpls)
	if err != nil {
		log.Fatalf("failed to list current templates")
	}

	customGroupName := "custom-group-1"

	createInput := &services.CreateGroupInput{
		GroupName:           customGroupName,
		TemplateID:          templates[0].ID,
		Capacity:            2,
		HealthCheckInterval: 300,
	}
	err = svc.Groups().Create(context.Background(), createInput)
	if err != nil {
		log.Fatalf("failed to create service group: %v", err)
	}

	fmt.Printf("Created Group: %s\n", customGroupName)

	fmt.Println("---")

	deleteInput := &services.DeleteGroupInput{
		Name: customGroupName,
	}
	err = svc.Groups().Delete(context.Background(), deleteInput)
	if err != nil {
		log.Fatalf("failed to delete service group: %v", err)
	}

	fmt.Printf("Delete Group: %s\n", customGroupName)

}
