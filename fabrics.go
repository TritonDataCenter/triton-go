package triton

import (
	"encoding/json"
	"fmt"
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

type CreateFabricVLANInput struct {
	Name        string `json:"name"`
	ID          int    `json:"vlan_id"`
	Description string `json:"description"`
}

func (client *FabricsClient) CreateFabricVLAN(input *CreateFabricVLANInput) (*FabricVLAN, error) {
	path := fmt.Sprintf("/%s/fabrics/default/vlans", client.accountName)
	respReader, err := client.executeRequest(http.MethodPost, path, input)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing CreateFabricVLAN request: {{err}}", err)
	}

	var result *FabricVLAN
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding CreateFabricVLAN response: {{err}}", err)
	}

	return result, nil
}

type UpdateFabricVLANInput struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (client *FabricsClient) UpdateFabricVLAN(input *UpdateFabricVLANInput) (*FabricVLAN, error) {
	path := fmt.Sprintf("/%s/fabrics/default/vlans/%d", client.accountName, input.ID)
	respReader, err := client.executeRequest(http.MethodPut, path, input)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing UpdateFabricVLAN request: {{err}}", err)
	}

	var result *FabricVLAN
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding UpdateFabricVLAN response: {{err}}", err)
	}

	return result, nil
}

type GetFabricVLANInput struct {
	ID int `json:"-"`
}

func (client *FabricsClient) GetFabricVLAN(input *GetFabricVLANInput) (*FabricVLAN, error) {
	path := fmt.Sprintf("/%s/fabrics/default/vlans/%d", client.accountName, input.ID)
	respReader, err := client.executeRequest(http.MethodGet, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetFabricVLAN request: {{err}}", err)
	}

	var result *FabricVLAN
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetFabricVLAN response: {{err}}", err)
	}

	return result, nil
}

type DeleteFabricVLANInput struct {
	ID int `json:"-"`
}

func (client *FabricsClient) DeleteFabricVLAN(input *DeleteFabricVLANInput) error {
	path := fmt.Sprintf("/%s/fabrics/default/vlans/%d", client.accountName, input.ID)
	respReader, err := client.executeRequest(http.MethodDelete, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteFabricVLAN request: {{err}}", err)
	}

	return nil
}
