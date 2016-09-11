package triton

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jen20/triton-go/authentication"
)

// Client represents a connection to the Triton API.
type Client struct {
	client      *retryablehttp.Client
	authorizer  []authentication.Signer
	endpoint    string
	accountName string
}

// NewClient is used to construct a Client in order to make API
// requests to the Triton API.
//
// At least one signer must be provided - example signers include
// authentication.PrivateKeySigner and authentication.SSHAgentSigner.
func NewClient(endpoint string, accountName string, signers ...authentication.Signer) (*Client, error) {
	return &Client{
		client:      retryablehttp.NewClient(),
		authorizer:  signers,
		endpoint:    strings.TrimSuffix(endpoint, "/"),
		accountName: accountName,
	}, nil
}

// Keys returns a client used for accessing functions pertaining to
// SSH key functionality in the Triton API.
func (client *Client) Keys() *KeysClient {
	return &KeysClient{client}
}

func (c *Client) formatURL(path string) string {
	return fmt.Sprintf("%s%s", c.endpoint, path)
}

func (c *Client) executeRequest(method, path string, body interface{}) (io.ReadCloser, error) {
	var requestBody io.ReadSeeker
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	req, err := retryablehttp.NewRequest(method, c.formatURL(path), requestBody)
	if err != nil {
		return nil, errwrap.Wrapf("Error constructing HTTP request: {{err}}", err)
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	authHeader, err := c.authorizer[0].Sign(dateHeader)
	if err != nil {
		return nil, errwrap.Wrapf("Error signing HTTP request: {{err}}", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Version", "8")
	req.Header.Set("User-Agent", "triton-go client API")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errwrap.Wrapf("Error executing HTTP request: {{err}}", err)
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, nil
	}

	tritonError := &TritonError{
		StatusCode: resp.StatusCode,
	}

	errorDecoder := json.NewDecoder(resp.Body)
	if err := errorDecoder.Decode(tritonError); err != nil {
		return nil, errwrap.Wrapf("Error decoding error resopnse: {{err}}", err)
	}
	return nil, tritonError
}
