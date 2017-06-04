package manta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jen20/manta-go/authentication"
)

// Client represents a connection to the Triton API.
type Client struct {
	client      *retryablehttp.Client
	authorizer  []authentication.Signer
	endpoint    string
	accountName string
	userAgent   string
}

type ClientOptions struct {
	Endpoint    string
	AccountName string
	UserAgent   string
	Signers     []authentication.Signer
}

// NewClient is used to construct a Client in order to make API
// requests to the Triton API.
//
// At least one signer must be provided - example signers include
// authentication.PrivateKeySigner and authentication.SSHAgentSigner.
func NewClient(options *ClientOptions) (*Client, error) {
	defaultRetryWaitMin := 1 * time.Second
	defaultRetryWaitMax := 5 * time.Minute
	defaultRetryMax := 32

	httpClient := &http.Client{
		Transport:     cleanhttp.DefaultTransport(),
		CheckRedirect: doNotFollowRedirects,
	}

	retryableClient := &retryablehttp.Client{
		HTTPClient:   httpClient,
		Logger:       log.New(os.Stderr, "", log.LstdFlags),
		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
	}

	client := &Client{
		client:      retryableClient,
		authorizer:  options.Signers,
		endpoint:    strings.TrimSuffix(options.Endpoint, "/"),
		accountName: options.AccountName,
	}

	if options.UserAgent == "" {
		client.userAgent = "Joyent manta-go Client SDK"
	} else {
		client.userAgent = options.UserAgent
	}

	return client, nil
}

func doNotFollowRedirects(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}

func (c *Client) formatURL(path string) string {
	return fmt.Sprintf("%s%s", c.endpoint, path)
}

func (c *Client) executeRequest(method, path string, query *url.Values, headers *http.Header, body interface{}) (io.ReadCloser, http.Header, error) {
	var requestBody io.ReadSeeker
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	req, err := retryablehttp.NewRequest(method, c.formatURL(path), requestBody)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error constructing HTTP request: {{err}}", err)
	}

	if body != nil && (headers == nil || headers.Get("Content-Type") == "") {
		req.Header.Set("Content-Type", "application/json")
	}
	if headers != nil {
		for key, values := range *headers {
			for _, value := range values {
				req.Header.Set(key, value)
			}
		}
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	authHeader, err := c.authorizer[0].Sign(dateHeader)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error signing HTTP request: {{err}}", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "manta-go client API")

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error executing HTTP request: {{err}}", err)
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, resp.Header, nil
	}

	mantaError := &MantaError{
		StatusCode: resp.StatusCode,
	}

	errorDecoder := json.NewDecoder(resp.Body)
	if err := errorDecoder.Decode(mantaError); err != nil {
		return nil, nil, errwrap.Wrapf("Error decoding error response: {{err}}", err)
	}
	return nil, nil, mantaError
}

func (c *Client) executeRequestNoEncode(method, path string, query *url.Values, headers *http.Header, body io.ReadSeeker) (io.ReadCloser, http.Header, error) {
	req, err := retryablehttp.NewRequest(method, c.formatURL(path), body)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error constructing HTTP request: {{err}}", err)
	}

	if headers != nil {
		for key, values := range *headers {
			for _, value := range values {
				req.Header.Set(key, value)
			}
		}
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	authHeader, err := c.authorizer[0].Sign(dateHeader)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error signing HTTP request: {{err}}", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "manta-go client API")

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error executing HTTP request: {{err}}", err)
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, resp.Header, nil
	}

	mantaError := &MantaError{
		StatusCode: resp.StatusCode,
	}

	errorDecoder := json.NewDecoder(resp.Body)
	if err := errorDecoder.Decode(mantaError); err != nil {
		return nil, nil, errwrap.Wrapf("Error decoding error response: {{err}}", err)
	}
	return nil, nil, mantaError
}
