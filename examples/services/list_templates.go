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

	s, err := services.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create new services client: %v", err)
	}

	listInput := &services.ListTemplatesInput{}
	templates, err := s.Templates().List(context.Background(), listInput)
	if err != nil {
		log.Fatalf("failed to list instance templates: %v", err)
	}
	for _, template := range templates {
		fmt.Printf("Template Name: %s\n", template.TemplateName)
		fmt.Printf("ID: %v\n", template.ID)
		fmt.Printf("AccountId: %v\n", template.AccountId)
		fmt.Printf("Package: %v\n", template.Package)
		fmt.Printf("ImageId: %v\n", template.ImageId)
		fmt.Printf("InstanceNamePrefix: %v\n", template.InstanceNamePrefix)
		fmt.Printf("FirewallEnabled: %v\n", template.FirewallEnabled)
		fmt.Printf("Networks: %v\n", template.Networks)
		fmt.Printf("UserData: %v\n", template.UserData)
		fmt.Printf("MetaData: %v\n", template.MetaData)
		fmt.Printf("Tags: %v\n", template.Tags)
		fmt.Println("")
	}

	fmt.Println("---")

	if tmpl := templates[0]; tmpl != nil {
		getInput := &services.GetTemplateInput{
			Name: tmpl.TemplateName,
		}
		template, err := s.Templates().Get(context.Background(), getInput)
		if err != nil {
			log.Fatalf("failed to get instance template: %v", err)
		}

		fmt.Printf("Got Template: %s\n", template.TemplateName)
	}

	fmt.Println("---")

	customTemplateName := "custom-template-1"

	createInput := &services.CreateTemplateInput{
		TemplateName:       customTemplateName,
		AccountId:          "joyent",
		Package:            "test-package",
		ImageId:            "49b22aec-0c8a-11e6-8807-a3eb4db576ba",
		InstanceNamePrefix: "sample-",
		FirewallEnabled:    false,
		Networks:           []string{"f7ed95d3-faaf-43ef-9346-15644403b963"},
		UserData:           "bash script here",
		MetaData:           map[string]string{"metadata": "test"},
		Tags:               map[string]string{"tag": "test"},
	}
	err = s.Templates().Create(context.Background(), createInput)
	if err != nil {
		log.Fatalf("failed to create instance template: %v", err)
	}

	fmt.Printf("Created Template: %s\n", customTemplateName)
	fmt.Println("---")

	deleteInput := &services.DeleteTemplateInput{
		Name: customTemplateName,
	}
	err = s.Templates().Delete(context.Background(), deleteInput)
	if err != nil {
		log.Fatalf("failed to delete instance template: %v", err)
	}

	fmt.Printf("Delete Template: %s\n", customTemplateName)

}
