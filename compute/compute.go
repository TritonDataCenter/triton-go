package compute

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type Compute struct {
	client *client.Client
}

// Accounts returns a Compute client used for accessing functions pertaining to
// Account functionality in the Triton API.
func (c *Compute) Accounts() *AccountsClient {
	return &AccountsClient{c}
}

// Config returns a c used for accessing functions pertaining
// to Config functionality in the Triton API.
func (c *Compute) Config() *ConfigClient {
	return &ConfigClient{c}
}

// Datacenters returns a Compute client used for accessing functions pertaining
// to DataCenter functionality in the Triton API.
func (c *Compute) Datacenters() *DataCentersClient {
	return &DataCentersClient{c}
}

// Fabrics returns a Compute client used for accessing functions pertaining to
// Fabric functionality in the Triton API.
func (c *Compute) Fabrics() *FabricsClient {
	return &FabricsClient{c}
}

// Firewall returns a Compute client used for accessing functions pertaining to
// firewall functionality in the Triton API.
func (c *Compute) Firewall() *FirewallClient {
	return &FirewallClient{c}
}

// Images returns a Compute client used for accessing functions pertaining to
// Images functionality in the Triton API.
func (c *Compute) Images() *ImagesClient {
	return &ImagesClient{c}
}

// Keys returns a Compute client used for accessing functions pertaining to SSH
// key functionality in the Triton API.
func (c *Compute) Keys() *KeysClient {
	return &KeysClient{c}
}

// Machines returns a Compute client used for accessing functions pertaining to
// machine functionality in the Triton API.
func (c *Compute) Machines() *MachinesClient {
	return &MachinesClient{c}
}

// Networks returns a Compute client used for accessing functions pertaining to
// Network functionality in the Triton API.
func (c *Compute) Networks() *NetworksClient {
	return &NetworksClient{c}
}

// Packages returns a Compute client used for accessing functions pertaining to
// Packages functionality in the Triton API.
func (c *Compute) Packages() *PackagesClient {
	return &PackagesClient{c}
}

// Roles returns a Compute client used for accessing functions pertaining to
// Role functionality in the Triton API.
func (c *Compute) Roles() *RolesClient {
	return &RolesClient{c}
}

// Services returns a Compute client used for accessing functions pertaining to
// Services functionality in the Triton API.
func (c *Compute) Services() *ServicesClient {
	return &ServicesClient{c}
}

// -----------------------------------------------------------------------------

func (c *Compute) executeRequestURIParams(ctx context.Context, method, path string, body interface{}, query *url.Values) (io.ReadCloser, error) {
	var requestBody io.ReadSeeker
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	endpoint := c.client.APIURL
	endpoint.Path = path
	if query != nil {
		endpoint.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, endpoint.String(), requestBody)
	if err != nil {
		return nil, errwrap.Wrapf("Error constructing HTTP request: {{err}}", err)
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	// NewClient ensures there's always an authorizer (unless this is called
	// outside that constructor).
	authHeader, err := c.client.Authorizers[0].Sign(dateHeader)
	if err != nil {
		return nil, errwrap.Wrapf("Error signing HTTP request: {{err}}", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Version", "8")
	req.Header.Set("User-Agent", "triton-go Client API")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errwrap.Wrapf("Error executing HTTP request: {{err}}", err)
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, nil
	}

	return nil, c.client.DecodeError(resp.StatusCode, resp.Body)
}

func (c *Compute) executeRequest(ctx context.Context, method, path string, body interface{}) (io.ReadCloser, error) {
	return c.executeRequestURIParams(ctx, method, path, body, nil)
}

func (c *Compute) executeRequestRaw(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var requestBody io.ReadSeeker
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	endpoint := c.client.APIURL
	endpoint.Path = path

	req, err := http.NewRequest(method, endpoint.String(), requestBody)
	if err != nil {
		return nil, errwrap.Wrapf("Error constructing HTTP request: {{err}}", err)
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	// NewClient ensures there's always an authorizer (unless this is called
	// outside that constructor).
	authHeader, err := c.client.Authorizers[0].Sign(dateHeader)
	if err != nil {
		return nil, errwrap.Wrapf("Error signing HTTP request: {{err}}", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Version", "8")
	req.Header.Set("User-Agent", "triton-go c API")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errwrap.Wrapf("Error executing HTTP request: {{err}}", err)
	}

	return resp, nil
}
