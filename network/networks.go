package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type NetworksClient struct {
	client *client.Client
}

type Network struct {
	Id                  string            `json:"id"`
	Name                string            `json:"name"`
	Public              bool              `json:"public"`
	Fabric              bool              `json:"fabric"`
	Description         string            `json:"description"`
	Subnet              string            `json:"subnet"`
	ProvisioningStartIP string            `json:"provision_start_ip"`
	ProvisioningEndIP   string            `json:"provision_end_ip"`
	Gateway             string            `json:"gateway"`
	Resolvers           []string          `json:"resolvers"`
	Routes              map[string]string `json:"routes"`
	InternetNAT         bool              `json:"internet_nat"`
}

type ListNetworksInput struct{}

func (c *NetworksClient) ListNetworks(ctx context.Context, _ *ListNetworksInput) ([]*Network, error) {
	path := fmt.Sprintf("/%s/networks", c.client.AccountName)
	respReader, err := c.client.ExecuteRequest(ctx, http.MethodGet, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListNetworks request: {{err}}", err)
	}

	var result []*Network
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListNetworks response: {{err}}", err)
	}

	return result, nil
}

type GetNetworkInput struct {
	ID string
}

func (c *NetworksClient) GetNetwork(ctx context.Context, input *GetNetworkInput) (*Network, error) {
	path := fmt.Sprintf("/%s/networks/%s", c.client.AccountName, input.ID)
	respReader, err := c.client.ExecuteRequest(ctx, http.MethodGet, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetNetwork request: {{err}}", err)
	}

	var result *Network
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetNetwork response: {{err}}", err)
	}

	return result, nil
}
