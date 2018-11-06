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
	"strconv"

	"encoding/pem"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/storage"
)

func main() {
	keyID := os.Getenv("TRITON_KEY_ID")
	accountName := os.Getenv("TRITON_ACCOUNT")
	keyMaterial := os.Getenv("TRITON_KEY_MATERIAL")
	userName := os.Getenv("TRITON_USER")
	fileName := "foo.txt"
	localPath := "/tmp/" + fileName
	mantaPath := "/stor/bar/baz/" + fileName

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
		MantaURL:    os.Getenv("MANTA_URL"),
		AccountName: accountName,
		Username:    userName,
		Signers:     []authentication.Signer{signer},
	}

	client, err := storage.NewClient(config)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}

	mpuBody := storage.CreateMpuBody{
		ObjectPath: mantaPath,
	}

	fooFile := []byte("this is only a test\n")
	err = ioutil.WriteFile(localPath, fooFile, 0644)
	if err != nil {
		log.Fatalf("Failed to write temporary upload file " + localPath)
	}

	createMpuInput := &storage.CreateMpuInput{
		DurabilityLevel: 2,
		Body:            mpuBody,
		ForceInsert:     true,
	}

	// Create a multipart upload to use for further testing
	fmt.Println("*** Creating new multipart upload ***")
	response := &storage.CreateMpuOutput{}
	response, err = client.Objects().CreateMultipartUpload(context.Background(), createMpuInput)
	if err != nil {
		log.Fatalf("storage.Objects.CreateMpu: %v", err)
	}
	fmt.Printf("Response Body\nid: %s\npartsDirectory: %s\n", response.Id, response.PartsDirectory)
	fmt.Println("Successfully created MPU for " + localPath)

	reader, err := os.Open(localPath)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer reader.Close()

	uploadPartInput := &storage.UploadPartInput{
		Id:           response.Id,
		PartNum:      0,
		ObjectReader: reader,
	}

	fmt.Println("ObjectDirectorPath for UploadPartInput: " + response.PartsDirectory)

	// Upload a single part
	fmt.Println("\n*** Upload a single part to the previous multipart upload ***")
	response2 := &storage.UploadPartOutput{}
	response2, err = client.Objects().UploadPart(context.Background(), uploadPartInput)
	if err != nil {
		log.Fatalf("storage.Objects.UploadPart: %v", err)
	}
	fmt.Println("Successfully uploaded " + fileName + " part 0!")

	var parts []string
	fmt.Printf("Part: %s\n", response2.Part)
	parts = append(parts, response2.Part)
	commitBody := storage.CommitMpuBody{
		Parts: parts,
	}

	commitMpuInput := &storage.CommitMpuInput{
		Id:   response.Id,
		Body: commitBody,
	}

	// List parts
	fmt.Println("\n*** List the parts of the current multipart upload ***")
	listMpuInput := &storage.ListMpuPartsInput{
		Id: response.Id,
	}
	listPartsOutput, err := client.Objects().ListMultipartUploadParts(context.Background(), listMpuInput)
	if err != nil {
		log.Fatalf("storage.Objects.ListMultipartUploadParts: %v", err)
	}
	for _, value := range listPartsOutput.Parts {
		fmt.Println("Etag: " + value.ETag + " PartNumber: " + strconv.Itoa(value.PartNumber) + " Size: " + strconv.FormatInt(value.Size, 10))
	}
	fmt.Println("Successfully listed MPU parts!")

	// Commit completed multipart upload
	fmt.Println("\n*** Commit the completed multipart upload ***")
	err = client.Objects().CommitMultipartUpload(context.Background(), commitMpuInput)
	if err != nil {
		log.Fatalf("storage.Objects.CommitMultipartUpload: %v", err)
	}
	fmt.Println("Successfully committed " + response.Id + "!")

	getMpuInput := &storage.GetMpuInput{
		PartsDirectoryPath: response.PartsDirectory,
	}

	// Get the status of the completed multipart upload
	fmt.Println("\n*** Get the status of the multipart upload ***")
	response3 := &storage.GetMpuOutput{}
	response3, err = client.Objects().GetMultipartUpload(context.Background(), getMpuInput)
	if err != nil {
		log.Fatalf("storage.Objects.GetMultipartUpload: %v", err)
	}
	fmt.Println("Successful get of " + response3.Id + " for targetObject: " + response3.TargetObject)

	err = os.Remove(localPath)
	if err != nil {
		log.Fatalf("os.Remove: %v", err)
	}

	// Create a new multipart upload just to test abort
	fmt.Println("\n*** Create a throwaway multipart upload ***")
	response, err = client.Objects().CreateMultipartUpload(context.Background(), createMpuInput)
	if err != nil {
		log.Fatalf("storage.Objects.CreateMpu: %v", err)
	}
	fmt.Printf("Response Body\nid: %s\npartsDirectory: %s\n", response.Id, response.PartsDirectory)
	fmt.Println("Successfully created MPU for " + localPath)

	abortMpuInput := &storage.AbortMpuInput{
		PartsDirectoryPath: response.PartsDirectory,
	}

	// Abort multipart upload
	fmt.Println("\n*** Abort the previous multipart upload ***")
	err = client.Objects().AbortMultipartUpload(context.Background(), abortMpuInput)
	if err != nil {
		log.Fatalf("storage.Objects.AbortMultipartUpload: %v", err)
	}

	fmt.Println("Successfully aborted MPU")
}
