package storage

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/joyent/triton-go/authentication"
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
