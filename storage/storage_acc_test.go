//
// Copyright 2026 Edgecast Cloud LLC.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package storage_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	triton "github.com/TritonDataCenter/triton-go/v2"
	"github.com/TritonDataCenter/triton-go/v2/authentication"
	terrors "github.com/TritonDataCenter/triton-go/v2/errors"
	"github.com/TritonDataCenter/triton-go/v2/storage"
)

func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func buildStorageClient(t *testing.T) *storage.StorageClient {
	t.Helper()

	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST' set")
	}

	mantaURL := triton.GetEnv("MANTA_URL")
	if mantaURL == "" {
		t.Skip("TRITON_MANTA_URL or SDC_MANTA_URL must be set to run storage acceptance tests")
	}

	tritonURL := triton.GetEnv("URL")
	tritonAccount := triton.GetEnv("ACCOUNT")
	tritonKeyID := triton.GetEnv("KEY_ID")
	tritonKeyMaterial := triton.GetEnv("KEY_MATERIAL")
	userName := triton.GetEnv("USER")

	if tritonAccount == "" {
		t.Fatal("TRITON_ACCOUNT must be set")
	}
	if tritonKeyID == "" {
		t.Fatal("TRITON_KEY_ID must be set")
	}

	var signer authentication.Signer
	var err error
	if tritonKeyMaterial != "" {
		signer, err = authentication.NewPrivateKeySigner(authentication.PrivateKeySignerInput{
			KeyID:              tritonKeyID,
			PrivateKeyMaterial: []byte(tritonKeyMaterial),
			AccountName:        tritonAccount,
			Username:           userName,
		})
	} else {
		signer, err = authentication.NewSSHAgentSigner(authentication.SSHAgentSignerInput{
			KeyID:       tritonKeyID,
			AccountName: tritonAccount,
			Username:    userName,
		})
	}
	if err != nil {
		t.Fatalf("Error creating signer: %s", err)
	}

	client, err := storage.NewClient(&triton.ClientConfig{
		TritonURL:   tritonURL,
		MantaURL:    mantaURL,
		AccountName: tritonAccount,
		Username:    userName,
		Signers:     []authentication.Signer{signer},
	})
	if err != nil {
		t.Fatalf("Error creating storage client: %s", err)
	}

	return client
}

// TestAccStorageLifecycle exercises the full Manta object storage workflow:
// directory creation, object upload, listing, retrieval, deletion semantics,
// and recursive cleanup.
func TestAccStorageLifecycle(t *testing.T) {
	client := buildStorageClient(t)
	ctx := context.Background()

	testID := randomHex(16)
	basePath := "/stor/" + testID
	nestedPath := basePath + "/" + testID
	objectContent := testID

	// Always clean up, even on failure.
	t.Cleanup(func() {
		_ = client.Dir().Delete(ctx, &storage.DeleteDirectoryInput{
			DirectoryName: basePath,
			ForceDelete:   true,
		})
	})

	// Step 1: Create nested directory structure.
	t.Log("Step 1: Creating directories")
	if err := client.Dir().Put(ctx, &storage.PutDirectoryInput{
		DirectoryName: basePath,
	}); err != nil {
		t.Fatalf("Dir().Put(%s): %v", basePath, err)
	}
	if err := client.Dir().Put(ctx, &storage.PutDirectoryInput{
		DirectoryName: nestedPath,
	}); err != nil {
		t.Fatalf("Dir().Put(%s): %v", nestedPath, err)
	}

	// Step 2: Upload objects.
	t.Log("Step 2: Uploading objects")
	if err := client.Objects().Put(ctx, &storage.PutObjectInput{
		ObjectPath:  nestedPath + "/uuid.txt",
		ContentType: "text/plain",
		ObjectReader: strings.NewReader(objectContent),
	}); err != nil {
		t.Fatalf("Objects().Put(%s/uuid.txt): %v", nestedPath, err)
	}
	if err := client.Objects().Put(ctx, &storage.PutObjectInput{
		ObjectPath:  basePath + "/uuid.txt",
		ContentType: "text/plain",
		ObjectReader: strings.NewReader(objectContent),
	}); err != nil {
		t.Fatalf("Objects().Put(%s/uuid.txt): %v", basePath, err)
	}

	// Step 3: List top-level directory.
	t.Log("Step 3: Listing top-level directory")
	topList, err := client.Dir().List(ctx, &storage.ListDirectoryInput{
		DirectoryName: basePath,
	})
	if err != nil {
		t.Fatalf("Dir().List(%s): %v", basePath, err)
	}
	foundDir := false
	foundObj := false
	for _, entry := range topList.Entries {
		if entry.Name == testID && entry.Type == "directory" {
			foundDir = true
		}
		if entry.Name == "uuid.txt" && entry.Type == "object" {
			foundObj = true
		}
	}
	if !foundDir {
		t.Errorf("Expected directory entry %q in listing of %s, got: %+v",
			testID, basePath, topList.Entries)
	}
	if !foundObj {
		t.Errorf("Expected object entry \"uuid.txt\" in listing of %s, got: %+v",
			basePath, topList.Entries)
	}

	// Step 4: List nested directory.
	t.Log("Step 4: Listing nested directory")
	nestedList, err := client.Dir().List(ctx, &storage.ListDirectoryInput{
		DirectoryName: nestedPath,
	})
	if err != nil {
		t.Fatalf("Dir().List(%s): %v", nestedPath, err)
	}
	foundNestedObj := false
	for _, entry := range nestedList.Entries {
		if entry.Name == "uuid.txt" && entry.Type == "object" {
			foundNestedObj = true
		}
	}
	if !foundNestedObj {
		t.Errorf("Expected object entry \"uuid.txt\" in listing of %s, got: %+v",
			nestedPath, nestedList.Entries)
	}

	// Step 5: Get and verify object contents.
	t.Log("Step 5: Verifying object contents")
	for _, objPath := range []string{
		basePath + "/uuid.txt",
		nestedPath + "/uuid.txt",
	} {
		output, err := client.Objects().Get(ctx, &storage.GetObjectInput{
			ObjectPath: objPath,
		})
		if err != nil {
			t.Fatalf("Objects().Get(%s): %v", objPath, err)
		}
		body, err := ioutil.ReadAll(output.ObjectReader)
		output.ObjectReader.Close()
		if err != nil {
			t.Fatalf("reading body of %s: %v", objPath, err)
		}
		if string(body) != objectContent {
			t.Errorf("Objects().Get(%s): expected body %q, got %q",
				objPath, objectContent, string(body))
		}
	}

	// Step 6: Delete non-empty nested directory (should fail).
	t.Log("Step 6: Attempting to delete non-empty directory")
	err = client.Dir().Delete(ctx, &storage.DeleteDirectoryInput{
		DirectoryName: nestedPath,
		ForceDelete:   false,
	})
	if err == nil {
		t.Fatal("Expected error deleting non-empty directory, got nil")
	}
	if !terrors.IsDirectoryNotEmptyError(err) {
		t.Fatalf("Expected DirectoryNotEmpty error, got: %v", err)
	}

	// Step 7: Delete nested object.
	t.Log("Step 7: Deleting nested object")
	if err := client.Objects().Delete(ctx, &storage.DeleteObjectInput{
		ObjectPath: nestedPath + "/uuid.txt",
	}); err != nil {
		t.Fatalf("Objects().Delete(%s/uuid.txt): %v", nestedPath, err)
	}

	// Step 8: Delete now-empty nested directory (should succeed).
	t.Log("Step 8: Deleting empty nested directory")
	if err := client.Dir().Delete(ctx, &storage.DeleteDirectoryInput{
		DirectoryName: nestedPath,
		ForceDelete:   false,
	}); err != nil {
		t.Fatalf("Dir().Delete(%s): %v", nestedPath, err)
	}

	// Step 9: Verify nested directory is gone.
	t.Log("Step 9: Verifying nested directory is gone")
	_, err = client.Dir().List(ctx, &storage.ListDirectoryInput{
		DirectoryName: nestedPath,
	})
	if err == nil {
		t.Fatal("Expected error listing deleted directory, got nil")
	}
	if !terrors.IsResourceNotFoundError(err) && !terrors.IsDirectoryDoesNotExistError(err) &&
		!terrors.IsStatusNotFoundCode(err) {
		t.Fatalf("Expected ResourceNotFound or DirectoryDoesNotExist error, got: %v", err)
	}

	// Step 10: Recursive delete of remaining tree.
	t.Log("Step 10: Recursive delete of remaining tree")
	if err := client.Dir().Delete(ctx, &storage.DeleteDirectoryInput{
		DirectoryName: basePath,
		ForceDelete:   true,
	}); err != nil {
		t.Fatalf("Dir().Delete(%s, ForceDelete=true): %v", basePath, err)
	}

	// Step 11: Verify top-level directory is gone.
	t.Log("Step 11: Verifying top-level directory is gone")
	_, err = client.Dir().List(ctx, &storage.ListDirectoryInput{
		DirectoryName: basePath,
	})
	if err == nil {
		t.Fatal("Expected error listing deleted directory, got nil")
	}
	if !terrors.IsResourceNotFoundError(err) && !terrors.IsDirectoryDoesNotExistError(err) &&
		!terrors.IsStatusNotFoundCode(err) {
		fmt.Printf("Note: got error type %T: %v\n", err, err)
	}
}
