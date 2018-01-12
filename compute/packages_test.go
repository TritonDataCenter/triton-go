//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package compute_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/testutils"
)

var (
	fakePackageId = "7b17343c-94af-6266-e0e8-893a3b9993d0"
)

func TestListPackages(t *testing.T) {
	computeClient := MockIdentityClient()

	do := func(ctx context.Context, cc *compute.ComputeClient) ([]*compute.Package, error) {
		defer testutils.DeactivateClient()

		packages, err := cc.Packages().List(ctx, &compute.ListPackagesInput{})
		if err != nil {
			return nil, err
		}
		return packages, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages"), listPackagesSuccess)

		resp, err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages"), listPackagesEmpty)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %s", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages"), listPackagesBadDecode)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %s", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages"), listPackagesError)

		resp, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "unable to list packages") {
			t.Errorf("expected error to equal testError: found %v", err)
		}
	})
}

func TestGetPackage(t *testing.T) {
	computeClient := MockIdentityClient()

	do := func(ctx context.Context, cc *compute.ComputeClient) (*compute.Package, error) {
		defer testutils.DeactivateClient()

		pkg, err := cc.Packages().Get(ctx, &compute.GetPackageInput{
			ID: fakePackageId,
		})
		if err != nil {
			return nil, err
		}
		return pkg, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages", fakePackageId), getPackageSuccess)

		resp, err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages", fakePackageId), getPackageEmpty)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %s", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages", fakePackageId), getPackageBadDecode)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %s", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages", "not-a-real-package-id"), getPackageError)

		resp, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "unable to get package") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func listPackagesSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`[
	{
    "id": "7b17343c-94af-6266-e0e8-893a3b9993d0",
    "name": "sdc_128",
    "memory": 128,
    "disk": 12288,
    "swap": 256,
    "vcpus": 1,
    "lwps": 1000,
    "default": false,
    "version": "1.0.0"
  }]
`)

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listPackagesEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func listPackagesBadDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`[{
    "id": "7b17343c-94af-6266-e0e8-893a3b9993d0",
    "name": "sdc_128",
    "memory": 128,
    "disk": 12288,
    "swap": 256,
    "vcpus": 1,
    "lwps": 1000,
    "default": false,
    "version": "1.0.0",
  }]`)

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listPackagesError(req *http.Request) (*http.Response, error) {
	return nil, errors.New("unable to list packages")
}

func getPackageSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
  "id": "7b17343c-94af-6266-e0e8-893a3b9993d0",
  "name": "sdc_128",
  "memory": 128,
  "disk": 12288,
  "swap": 256,
  "vcpus": 1,
  "lwps": 1000,
  "default": false,
  "version": "1.0.0"
}
`)

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getPackageBadDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
  "id": "7b17343c-94af-6266-e0e8-893a3b9993d0",
  "name": "sdc_128",
  "memory": 128,
  "disk": 12288,
  "swap": 256,
  "vcpus": 1,
  "lwps": 1000,
  "default": false,
  "version": "1.0.0",
}`)

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getPackageEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func getPackageError(req *http.Request) (*http.Response, error) {
	return nil, errors.New("unable to get package")
}
