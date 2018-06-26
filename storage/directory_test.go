//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package storage_test

import (
	"context"
	"net/http"
	"net/url"

	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/joyent/triton-go/storage"
	"github.com/joyent/triton-go/testutils"
	"github.com/pkg/errors"
)

var (
	errListDir     = errors.New("unable to list dir")
	dirPath        = "/stor/foobar.json"
	brokenDirPath  = "/missingfolder/foo.json"
	dirListingFull = `{"name":"subdirectory0","type":"directory","mtime":"2018-01-01T00:00:00.000Z"}
{"name":"subdirectory1","type":"directory","mtime":"2018-01-01T00:00:00.000Z"}
`
	dirListingPartial = `{"name":"subdirectory1","type":"directory","mtime":"2018-01-01T00:00:00.000Z"}
`
	dirLastEntry = "subdirectory1"
)

func TestList(t *testing.T) {
	storageClient := &storage.StorageClient{
		Client: testutils.NewMockClient(testutils.MockClientInput{
			AccountName: accountUrl,
		}),
	}

	do := func(ctx context.Context, sc *storage.StorageClient, marker string) (*storage.ListDirectoryOutput, error) {
		defer testutils.DeactivateClient()

		listInput := &storage.ListDirectoryInput{
			DirectoryName: dirPath,
			// Limit         uint64
			// Marker        string
		}

		if marker != "" {
			listInput.Marker = marker
		}

		return sc.Dir().List(ctx, listInput)
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, dirPath), listDirSuccess)
		expectedResultSetSize := 2

		output, err := do(context.Background(), storageClient, "")
		if err != nil {
			t.Fatal(err)
		}

		if uint64(len(output.Entries)) != output.ResultSetSize {
			t.Fatalf(
				"expected entries length to equal result-set-size: len(output.Entries) = %v, output.ResultSetSize = %v",
				len(output.Entries),
				output.ResultSetSize)
		}

		if len(output.Entries) != expectedResultSetSize {
			t.Fatalf(
				"expected entries length for simple listing to equal %v, found %v",
				expectedResultSetSize,
				len(output.Entries))
		}

	})

	t.Run("successfulWithMarker", func(t *testing.T) {
		v := url.Values{}
		v.Set("marker", dirLastEntry)

		testutils.RegisterResponder(
			"GET",
			path.Join("/", accountUrl, dirPath)+"?"+v.Encode(),
			listDirSuccess)

		expectedPartialResultSetSize := 1

		output, err := do(context.Background(), storageClient, dirLastEntry)
		if err != nil {
			t.Fatal(err)
		}

		if uint64(len(output.Entries)) != output.ResultSetSize {
			t.Fatalf(
				"expected entries length to equal result-set-size: len(output.Entries) = %v, output.ResultSetSize = %v",
				len(output.Entries),
				output.ResultSetSize)
		}

		if len(output.Entries) != expectedPartialResultSetSize {
			t.Fatalf(
				"expected entries length for partial listing to equal %v, found %v",
				expectedPartialResultSetSize,
				len(output.Entries))
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, brokenDirPath), listDirError)

		output, err := do(context.Background(), storageClient, "")
		if err == nil {
			t.Fatal("expected non-nil error, but err was nil")
		}

		if !strings.Contains(err.Error(), "unable to list dir") {
			t.Errorf("expected error to equal testError: found %v", err)
		}

		if output != nil {
			t.Fatalf("expected nil output: found %v", output)
		}
	})
}

func listDirSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/x-json-stream")

	var body *strings.Reader

	if req.URL.Query().Get("marker") == "subdirectory1" {
		header.Add("Result-Set-Size", "1")
		body = strings.NewReader(dirListingPartial)
	} else {
		header.Add("Result-Set-Size", "2")
		body = strings.NewReader(dirListingFull)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listDirError(req *http.Request) (*http.Response, error) {
	return nil, listDirErrorType
}
