//
// Copyright 2020 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package network_test

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/joyent/triton-go/network"
	"github.com/joyent/triton-go/testutils"
	"github.com/pkg/errors"
)

var (
	fakeNetworkID        = "daeb93a2-532e-4bd4-8788-b6b30f10ac17"
	getNetworkErrorType  = errors.New("unable to get network")
	listNetworkErrorType = errors.New("unable to list networks")
	testNetworkID        = ""
)

// Note that this is specific to Joyent Public Cloud and will not pass on
// private installations of Triton.
func TestAccNetworks_List(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepNetworkClient{
				StateBagKey: "networks",
				CallFunc: func(state testutils.TritonStateBag, c *network.NetworkClient) (interface{}, error) {
					ctx := context.Background()
					input := &network.ListInput{}

					networks, err := c.List(ctx, input)
					if err != nil {
						log.Fatalf("Networks.List(): %v", err)
					}

					if len(networks) == 0 {
						t.Fatal("No networks returned from network list")
					}

					testNetworkID = networks[0].Id

					return networks, nil
				},
			},
		},
	})
}

func TestAccNetworks_Get(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepNetworkClient{
				StateBagKey: "networks",
				CallFunc: func(state testutils.TritonStateBag, c *network.NetworkClient) (interface{}, error) {
					if testNetworkID == "" {
						t.Skip("No network id")
					}

					ctx := context.Background()
					input := &network.GetInput{
						ID: testNetworkID,
					}

					network, err := c.Get(ctx, input)
					if err != nil {
						log.Fatalf("Networks.List(): %v", err)
					}

					if network.Id != testNetworkID {
						t.Fatalf("Expected package id %s, got %s",
							testNetworkID, network.Id)
					}

					return nil, nil
				},
			},
		},
	})
}

func TestListNetworks(t *testing.T) {
	networkClient := MockNetworkClient()

	do := func(ctx context.Context, nc *network.NetworkClient) ([]*network.Network, error) {
		defer testutils.DeactivateClient()

		networks, err := nc.List(ctx, &network.ListInput{})
		if err != nil {
			return nil, err
		}
		return networks, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "networks"), listNetworksSuccess)

		resp, err := do(context.Background(), networkClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "networks"), listNetworksEmpty)

		_, err := do(context.Background(), networkClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %v", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "networks"), listNetworksBadeDecode)

		_, err := do(context.Background(), networkClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "networks"), listNetworksError)

		resp, err := do(context.Background(), networkClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "unable to list networks") {
			t.Errorf("expected error to equal testError: found %v", err)
		}
	})
}

func TestGetNetwork(t *testing.T) {
	networkClient := MockNetworkClient()

	do := func(ctx context.Context, nc *network.NetworkClient) (*network.Network, error) {
		defer testutils.DeactivateClient()

		network, err := nc.Get(ctx, &network.GetInput{
			ID: fakeNetworkID,
		})
		if err != nil {
			return nil, err
		}
		return network, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "networks", fakeNetworkID), getNetworkSuccess)

		resp, err := do(context.Background(), networkClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "networks", fakeNetworkID), getNetworkEmpty)

		_, err := do(context.Background(), networkClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %v", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "networks", fakeNetworkID), getNetworkBadeDecode)

		_, err := do(context.Background(), networkClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountURL, "networks"), getNetworkError)

		resp, err := do(context.Background(), networkClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "unable to get network") {
			t.Errorf("expected error to equal testError: found %v", err)
		}
	})
}

func getNetworkSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
  "id": "daeb93a2-532e-4bd4-8788-b6b30f10ac17",
  "name": "external",
  "public": true
}
`)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getNetworkError(req *http.Request) (*http.Response, error) {
	return nil, getNetworkErrorType
}

func getNetworkBadeDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
  "id": "daeb93a2-532e-4bd4-8788-b6b30f10ac17",
  "name": "external",
  "public": true,
}`)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getNetworkEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func listNetworksEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func listNetworksSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`[
	{
    "id": "daeb93a2-532e-4bd4-8788-b6b30f10ac17",
    "name": "external",
    "public": true
  }
]`)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listNetworksBadeDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{[
	{
    "id": "daeb93a2-532e-4bd4-8788-b6b30f10ac17",
    "name": "external",
    "public": true
  }
]}`)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listNetworksError(req *http.Request) (*http.Response, error) {
	return nil, listNetworkErrorType
}
