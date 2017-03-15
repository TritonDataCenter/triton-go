package triton

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/errwrap"
)

type FabricsClient struct {
	*Client
}

// Fabrics returns a client used for accessing functions pertaining to
// Fabric functionality in the Triton API.
func (c *Client) Fabrics() *FabricsClient {
	return &FabricsClient{c}
}

type FabricVLAN struct {
	Name        string
	ID          int `json:"vlan_id"`
	Description string
}

type ListFabricVLANsInput struct{}

func (client *FabricsClient) ListFabricVLANs(*ListFabricVLANsInput) ([]*FabricVLAN, error) {
	respReader, err := client.executeRequest(http.MethodGet, "/my/fabrics/default/vlans", nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListFabricVLANs request: {{err}}", err)
	}

	var result []*FabricVLAN
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListFabricVLANs response: {{err}}", err)
	}

	return result, nil
}
