package account

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type ConfigClient struct {
	client *client.Client
}

// Config represents configuration for your account.
type Config struct {
	// DefaultNetwork is the network that docker containers are provisioned on.
	DefaultNetwork string `json:"default_network"`
}

type GetConfigInput struct{}

// GetConfig outputs configuration for your account.
func (c *ConfigClient) Get(ctx context.Context, input *GetConfigInput) (*Config, error) {
	path := fmt.Sprintf("/%s/config", c.client.AccountName)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get account config")
	}

	var result *Config
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode get account config response")
	}

	return result, nil
}

type UpdateConfigInput struct {
	// DefaultNetwork is the network that docker containers are provisioned on.
	DefaultNetwork string `json:"default_network"`
}

// UpdateConfig updates configuration values for your account.
func (c *ConfigClient) Update(ctx context.Context, input *UpdateConfigInput) (*Config, error) {
	path := fmt.Sprintf("/%s/config", c.client.AccountName)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Body:   input,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to update account config")
	}

	var result *Config
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode update account config response")
	}

	return result, nil
}
