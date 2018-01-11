package compute

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

const pingEndpoint = "/--ping"

type CloudAPI struct {
	Versions []string `json:"versions"`
}

type PingOutput struct {
	Ping     string   `json:"ping"`
	CloudAPI CloudAPI `json:"cloudapi"`
}

// Ping sends a request to the '/--ping' endpoint and returns a `pong` as well
// as a list of API version numbers your instance of CloudAPI is presenting.
func (c *ComputeClient) Ping(ctx context.Context) (*PingOutput, error) {
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   pingEndpoint,
	}
	response, err := c.Client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return nil, errors.Wrap(err, "unable to ping")
	}
	if response == nil {
		return nil, errors.Wrap(err, "unable to ping")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
		return nil, &client.TritonError{
			StatusCode: response.StatusCode,
			Code:       "ResourceNotFound",
		}
	}

	var result *PingOutput
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&result); err != nil {
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode ping response")
		}
	}

	return result, nil
}
