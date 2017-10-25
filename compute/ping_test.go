package compute_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/client"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/testutils"
)

// test borked 404 response (TritonError)
// test borked 410 response (TritonError)
// test bad JSON decode

func TestPing(t *testing.T) {
	computeClient := buildClient()

	do := func(ctx context.Context, pc *compute.ComputeClient) (*compute.PingOutput, error) {
		defer testutils.DeactivateClient()

		ping, err := pc.Ping(ctx)
		if err != nil {
			return nil, err
		}
		return ping, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", "/--ping", pingSuccessFunc)

		resp, err := do(context.Background(), computeClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp.Ping != "pong" {
			t.Errorf("ping was not pong: expected %s", resp.Ping)
		}

		versions := []string{"7.0.0", "7.1.0", "7.2.0", "7.3.0", "8.0.0"}
		if !reflect.DeepEqual(resp.CloudAPI.Versions, versions) {
			t.Errorf("ping did not contain CloudAPI versions: expected %s", versions)
		}
	})

	t.Run("EOF decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", "/--ping", pingEmptyFunc)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %s", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", "/--ping", pingErrorFunc)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		blankError := "Ping request has empty response"

		if err.Error() != blankError {
			t.Errorf("expected error to equal defaultError: found %s", err)
		}
	})

	t.Run("404", func(t *testing.T) {
		testutils.RegisterResponder("GET", "/--ping", ping404Func)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "ResourceNotFound") {
			t.Errorf("expected error to be a 404: found %s", err)
		}
	})

	t.Run("410", func(t *testing.T) {
		testutils.RegisterResponder("GET", "/--ping", ping410Func)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "ResourceNotFound") {
			t.Errorf("expected error to be a 410: found %s", err)
		}
	})

	t.Run("bad decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", "/--ping", pingDecodeFunc)

		_, err := do(context.Background(), computeClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %s", err)
		}
	})
}

var defaultError = errors.New("we got the funk")
var defaultHeader = http.Header{}

func pingSuccessFunc(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
	"ping": "pong",
	"cloudapi": {
		"versions": ["7.0.0", "7.1.0", "7.2.0", "7.3.0", "8.0.0"]
	}
}`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func pingEmptyFunc(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func ping404Func(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: 404,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func ping410Func(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: 410,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func pingErrorFunc(req *http.Request) (*http.Response, error) {
	return nil, defaultError
}

func pingDecodeFunc(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
(ham!(//
}`)

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func buildClient() *compute.ComputeClient {
	testSigner, _ := authentication.NewTestSigner()
	httpClient := &client.Client{
		Authorizers: []authentication.Signer{testSigner},
		HTTPClient: &http.Client{
			Transport: testutils.DefaultMockTransport,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
	return &compute.ComputeClient{
		Client: httpClient,
	}
}
