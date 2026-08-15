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
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/TritonDataCenter/triton-go/v2/storage"
	"github.com/TritonDataCenter/triton-go/v2/testutils"
)

// dirEntry is a helper for building mock directory listings.
type dirEntry struct {
	Name string
	Type string // "directory" or "object"
}

// makeDirListing returns a newline-delimited JSON stream matching the format
// that Manta returns for directory listings (application/x-json-stream).
func makeDirListing(entries []dirEntry) string {
	var lines []string
	for _, e := range entries {
		lines = append(lines,
			fmt.Sprintf(`{"name":%q,"type":%q,"mtime":"2025-01-01T00:00:00.000Z"}`,
				e.Name, e.Type))
	}
	return strings.Join(lines, "\n") + "\n"
}

// deleteTracker returns a responder that records each DELETE request path and
// replies with 204 No Content.
func deleteTracker(mu *sync.Mutex, deleted map[string]bool) testutils.Responder {
	return func(req *http.Request) (*http.Response, error) {
		mu.Lock()
		deleted[req.URL.Path] = true
		mu.Unlock()
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header:     http.Header{"Content-Type": {"application/json"}},
		}, nil
	}
}

// listResponder returns a responder that serves a directory listing.
func listResponder(entries []dirEntry) testutils.Responder {
	body := makeDirListing(entries)
	count := fmt.Sprintf("%d", len(entries))
	return func(req *http.Request) (*http.Response, error) {
		header := http.Header{}
		header.Set("Content-Type", "application/x-json-stream")
		header.Set("Result-Set-Size", count)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       ioutil.NopCloser(strings.NewReader(body)),
		}, nil
	}
}

// TestForceDeleteMultipleObjects verifies that ForceDelete removes every object
// in a directory, not just the first one. It also checks that the target
// directory itself is deleted after its contents are removed.
func TestForceDeleteMultipleObjects(t *testing.T) {
	storageClient := MockStorageClient()
	defer testutils.DeactivateClient()

	var mu sync.Mutex
	deleted := make(map[string]bool)
	track := deleteTracker(&mu, deleted)

	// Directory listing: two objects
	testutils.RegisterResponder("GET",
		path.Join("/", accountURL, "stor/mydir"),
		listResponder([]dirEntry{
			{"file1.txt", "object"},
			{"file2.txt", "object"},
		}))

	// DELETE responders for the two objects and the directory itself
	testutils.RegisterResponder("DELETE",
		path.Join("/", accountURL, "stor/mydir/file1.txt"), track)
	testutils.RegisterResponder("DELETE",
		path.Join("/", accountURL, "stor/mydir/file2.txt"), track)
	testutils.RegisterResponder("DELETE",
		path.Join("/", accountURL, "stor/mydir"), track)

	err := storageClient.Dir().Delete(context.Background(), &storage.DeleteDirectoryInput{
		DirectoryName: "/stor/mydir",
		ForceDelete:   true,
	})
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Both objects must have been deleted
	expect := []string{
		path.Join("/", accountURL, "stor/mydir/file1.txt"),
		path.Join("/", accountURL, "stor/mydir/file2.txt"),
		path.Join("/", accountURL, "stor/mydir"),
	}
	for _, p := range expect {
		if !deleted[p] {
			t.Errorf("expected DELETE for %q but it was not called", p)
		}
	}
}

// TestForceDeleteNestedDirectory verifies that ForceDelete recurses into
// subdirectories and deletes the entire tree.
func TestForceDeleteNestedDirectory(t *testing.T) {
	storageClient := MockStorageClient()
	defer testutils.DeactivateClient()

	var mu sync.Mutex
	deleted := make(map[string]bool)
	track := deleteTracker(&mu, deleted)

	// Top-level listing: one subdirectory and one object
	testutils.RegisterResponder("GET",
		path.Join("/", accountURL, "stor/mydir"),
		listResponder([]dirEntry{
			{"subdir", "directory"},
			{"file1.txt", "object"},
		}))

	// Subdirectory listing: one object
	testutils.RegisterResponder("GET",
		path.Join("/", accountURL, "stor/mydir/subdir"),
		listResponder([]dirEntry{
			{"file3.txt", "object"},
		}))

	// DELETE responders for all four paths
	testutils.RegisterResponder("DELETE",
		path.Join("/", accountURL, "stor/mydir/subdir/file3.txt"), track)
	testutils.RegisterResponder("DELETE",
		path.Join("/", accountURL, "stor/mydir/subdir"), track)
	testutils.RegisterResponder("DELETE",
		path.Join("/", accountURL, "stor/mydir/file1.txt"), track)
	testutils.RegisterResponder("DELETE",
		path.Join("/", accountURL, "stor/mydir"), track)

	err := storageClient.Dir().Delete(context.Background(), &storage.DeleteDirectoryInput{
		DirectoryName: "/stor/mydir",
		ForceDelete:   true,
	})
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	expect := []string{
		path.Join("/", accountURL, "stor/mydir/subdir/file3.txt"),
		path.Join("/", accountURL, "stor/mydir/subdir"),
		path.Join("/", accountURL, "stor/mydir/file1.txt"),
		path.Join("/", accountURL, "stor/mydir"),
	}
	for _, p := range expect {
		if !deleted[p] {
			t.Errorf("expected DELETE for %q but it was not called", p)
		}
	}
}

// TestDeleteDirectorySimple verifies that Delete with ForceDelete: false
// issues a single DELETE request without listing or recursing.
func TestDeleteDirectorySimple(t *testing.T) {
	storageClient := MockStorageClient()
	defer testutils.DeactivateClient()

	var mu sync.Mutex
	deleted := make(map[string]bool)
	track := deleteTracker(&mu, deleted)

	dirPath := path.Join("/", accountURL, "stor/emptydir")

	// If a GET (directory listing) is issued, the test should fail.
	testutils.RegisterResponder("GET", dirPath,
		func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected directory listing request; ForceDelete is false")
			return nil, nil
		})

	testutils.RegisterResponder("DELETE", dirPath, track)

	err := storageClient.Dir().Delete(context.Background(), &storage.DeleteDirectoryInput{
		DirectoryName: "/stor/emptydir",
		ForceDelete:   false,
	})
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if !deleted[dirPath] {
		t.Errorf("expected DELETE for %q but it was not called", dirPath)
	}
	if len(deleted) != 1 {
		t.Errorf("expected exactly 1 DELETE, got %d: %v", len(deleted), deleted)
	}
}
