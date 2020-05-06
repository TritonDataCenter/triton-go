//
// Copyright 2020 Joyent, Inc.
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
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/testutils"
)

var (
	fakePackageId   = "7b17343c-94af-6266-e0e8-893a3b9993d0"
	fakePackageName = "g4-test"
	testPackageID   = ""
	testPackageName = ""
)

func TestAccPackagesList(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepComputeClient{
				StateBagKey: "package",
				CallFunc: func(state testutils.TritonStateBag, c *compute.ComputeClient) (interface{}, error) {
					ctx := context.Background()
					input := &compute.ListPackagesInput{}
					pkgs, err := c.Packages().List(ctx, input)
					if err != nil {
						log.Fatalf("Packages.List failed: %v", err)
					}

					if len(pkgs) == 0 {
						t.Fatal("No packages returned from package list")
					}

					testPackageID = pkgs[0].ID
					testPackageName = pkgs[0].Name
					for _, pkg := range pkgs {
						if len(pkg.Disks) != 0 {
							if pkg.Brand != "bhyve" {
								t.Fatalf("Package with disks does not have brand \"bhyve\": %+v", pkg.Disks)
							}
							for _, dsk := range pkg.Disks {
								if dsk.OSDiskSize {
									if dsk.Remaining {
										t.Errorf("OS Disk should not have Remaining set: %+v", dsk)
									}
								} else if dsk.Remaining {
									if dsk.SizeInMiB != 0 {
										t.Errorf("Disk with Remaining size should not have SizeInMiB: %+v", dsk)
									}
								} else {
									if dsk.SizeInMiB < 0 {
										t.Errorf("Disk SizeInMiB must be greater than zero: %+v", dsk)
									}
								}
							}
						}
					}
					if testPackageID == "" {
						t.Fatalf("Package does not have an ID %+v", pkgs[0])
					}

					return pkgs[0], nil
				},
			},
		},
	})
}

func TestAccPackagesListByName(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepComputeClient{
				StateBagKey: "package",
				CallFunc: func(state testutils.TritonStateBag, c *compute.ComputeClient) (interface{}, error) {
					if testPackageID == "" {
						t.Skip("No package id")
					}
					ctx := context.Background()
					input := &compute.ListPackagesInput{
						Name: testPackageName,
					}
					pkgs, err := c.Packages().List(ctx, input)
					if err != nil {
						log.Fatalf("Packages.List failed: %v", err)
					}

					if len(pkgs) == 0 {
						t.Fatal("No packages returned from package list")
					}

					for _, foundPkg := range pkgs {
						if foundPkg.Name != testPackageName {
							t.Fatalf("Expected package name %s, got %s",
								testPackageName, foundPkg.Name)
						}
					}

					return pkgs[0], nil
				},
			},
		},
	})
}

func TestAccPackagesGet(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepComputeClient{
				StateBagKey: "package",
				CallFunc: func(state testutils.TritonStateBag, c *compute.ComputeClient) (interface{}, error) {
					if testPackageID == "" {
						t.Skip("No package id")
					}
					ctx := context.Background()
					input := &compute.GetPackageInput{
						ID: testPackageID,
					}
					foundPkg, err := c.Packages().Get(ctx, input)
					if err != nil {
						log.Fatalf("Packages.Get failed: %v", err)
					}

					if foundPkg.ID != testPackageID {
						t.Fatalf("Expected package id %s, got %s",
							testPackageID, foundPkg.ID)
					}

					if foundPkg.Name != testPackageName {
						t.Fatalf("Expected package name %s, got %s",
							testPackageName, foundPkg.Name)
					}

					return foundPkg, nil
				},
			},
		},
	})
}

func TestListPackages(t *testing.T) {
	computeClient := MockComputeClient()

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

	t.Run("filtered", func(t *testing.T) {
		v := url.Values{}
		v.Set("name", fakePackageName)

		filterURL := path.Join("/", accountURL, "packages") + "?" + v.Encode()
		testutils.RegisterResponder("GET", filterURL, listPackagesFiltered)
		defer testutils.DeactivateClient()

		ctx := context.Background()
		cc := computeClient
		packages, err := cc.Packages().List(ctx, &compute.ListPackagesInput{
			Name: fakePackageName,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(packages) != 1 {
			t.Fatalf("expected output but received empty body")
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
	computeClient := MockComputeClient()

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

	t.Run("bhyve_disks", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "packages", fakePackageId), getPackageWithDisks)

		resp, err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}

		if resp.Brand != "bhyve" {
			t.Fatalf("Expected package brand to be bhyve")
		}

		if !resp.FlexibleDisk {
			t.Fatalf("Expected package to have FlexibleDisk")
		}

		if resp.Disks == nil {
			t.Fatalf("Expected package to have disks but got nil")
		}

		if len(resp.Disks) != 3 {
			t.Fatalf("Expected package to have 3 disks, but got %d",
				len(resp.Disks))
		}

		// Test disks are as expected.

		disk0 := resp.Disks[0]
		if disk0.Size != nil {
			t.Fatalf("Expected disk[0].Size to be nil, got %v",
				disk0.Size)
		}
		if disk0.SizeInMiB != 0 {
			t.Fatalf("Expected disk[0].SizeInMB to be 0, got %d",
				disk0.SizeInMiB)
		}
		if disk0.Remaining != false {
			t.Fatalf("Expected disk[0].Remaining to be false")
		}
		if disk0.OSDiskSize != true {
			t.Fatalf("Expected disk[0].OSDiskSize to be true")
		}

		disk1 := resp.Disks[1]
		if disk1.Size.(float64) != 6144 {
			t.Fatalf("Expected disk[0].Size to be 6144, got %v",
				disk1.Size)
		}
		if disk1.SizeInMiB != 6144 {
			t.Fatalf("Expected disk[0].SizeInMB to be 6144, got %d",
				disk1.SizeInMiB)
		}
		if disk1.Remaining != false {
			t.Fatalf("Expected disk[0].Remaining to be false")
		}
		if disk1.OSDiskSize != false {
			t.Fatalf("Expected disk[0].OSDiskSize to be false")
		}

		disk2 := resp.Disks[2]
		if disk2.Size != "remaining" {
			t.Fatalf("Expected disk[0].Size to be remaining, got %v",
				disk2.Size)
		}
		if disk2.SizeInMiB != 0 {
			t.Fatalf("Expected disk[0].SizeInMB to be 0, got %d",
				disk2.SizeInMiB)
		}
		if disk2.Remaining != true {
			t.Fatalf("Expected disk[0].Remaining to be true")
		}
		if disk2.OSDiskSize != false {
			t.Fatalf("Expected disk[0].OSDiskSize to be false")
		}
	})
}

func listPackagesFiltered(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`[
	{
	"id": "7b17343c-94af-6266-e0e8-893a3b9993d0",
	"name": "g4-test",
	"memory": 1024,
	"disk": 25600,
	"swap": 4096,
	"vcpus": 0,
	"lwps": 4000,
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

func getPackageWithDisks(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
    "brand": "bhyve",
    "default": true,
    "disk": 24576,
    "id": "7b17343c-94af-6266-e0e8-893a3b9993d0",
    "lwps": 4000,
    "memory": 1024,
    "name": "sample-bhyve-three-disks",
    "swap": 4096,
    "vcpus": 1,
    "version": "1.0.0",
    "flexible_disk": true,
    "disks": [
        {},
        {
            "size": 6144
        },
        {
            "size": "remaining"
        }
    ]
}`)

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}
