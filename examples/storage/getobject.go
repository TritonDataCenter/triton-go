package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"encoding/pem"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/storage"
)

func main() {
	var (
		signer authentication.Signer
		err    error

		keyID       = os.Getenv("MANTA_KEY_ID")
		accountName = os.Getenv("MANTA_USER")
		keyMaterial = os.Getenv("MANTA_KEY_MATERIAL")
	)

	if keyMaterial == "" {
		signer, err = authentication.NewSSHAgentSigner(keyID, accountName)
		if err != nil {
			log.Fatalf("Error Creating SSH Agent Signer: %s", err.Error())
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
			log.Fatalf("Error Creating SSH Private Key Signer: %s", err.Error())
		}
	}

	config := &triton.ClientConfig{
		MantaURL:    os.Getenv("MANTA_URL"),
		AccountName: accountName,
		Signers:     []authentication.Signer{signer},
	}

	client, err := storage.NewClient(config)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	// You must have this in Manta in order for this file to work.
	// path := "/stor/books/dracula.txt"
	path := "/stor/books/dracula.txt"

	ctx := context.Background()
	info, err := client.Objects().GetInfo(ctx, &storage.GetInfoInput{
		ObjectPath: path,
	})
	if err != nil {
		fmt.Printf("Could not find '%s'\n", path)
		return
	}

	fmt.Println("--- HEAD ---")
	fmt.Printf("Content-Length: %d\n", info.ContentLength)
	fmt.Printf("Content-MD5: %s\n", info.ContentMD5)
	fmt.Printf("Content-Type: %s\n", info.ContentType)
	fmt.Printf("ETag: %s\n", info.ETag)
	fmt.Printf("Date-Modified: %s\n", info.LastModified.String())

	ctx = context.Background()
	isDir, err := client.Objects().IsDir(ctx, path)
	if err != nil {
		log.Fatalf("failed to detect directory '%s'\n", err)
		return
	}

	if isDir {
		fmt.Printf("'%s' is a directory\n", path)
	} else {
		fmt.Printf("'%s' is a file\n", path)
	}

	ctx = context.Background()
	obj, err := client.Objects().Get(ctx, &storage.GetObjectInput{
		ObjectPath: path,
	})
	if err != nil {
		log.Fatalf("failed to get '%s': %s", path, err)
	}

	body, err := ioutil.ReadAll(obj.ObjectReader)
	if err != nil {
		log.Fatalf("failed to read response body: %s", err)
	}
	defer obj.ObjectReader.Close()

	fmt.Println("--- GET ---")
	fmt.Printf("Content-Length: %d\n", obj.ContentLength)
	fmt.Printf("Content-MD5: %s\n", obj.ContentMD5)
	fmt.Printf("Content-Type: %s\n", obj.ContentType)
	fmt.Printf("ETag: %s\n", obj.ETag)
	fmt.Printf("Date-Modified: %s\n", obj.LastModified.String())
	fmt.Printf("Length: %d\n", len(body))
}
