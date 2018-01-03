package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type PoliciesClient struct {
	client *client.Client
}

type Policy struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Rules       []string `json:"rules"`
	Description string   `json:"description"`
}

type ListPoliciesInput struct{}

func (c *PoliciesClient) List(ctx context.Context, _ *ListPoliciesInput) ([]*Policy, error) {
	path := fmt.Sprintf("/%s/policies", c.client.AccountName)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListPolicies request: {{err}}", err)
	}

	var result []*Policy
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListPolicies response: {{err}}", err)
	}

	return result, nil
}

type GetPolicyInput struct {
	PolicyID string
}

func (c *PoliciesClient) Get(ctx context.Context, input *GetPolicyInput) (*Policy, error) {
	path := fmt.Sprintf("/%s/policies/%s", c.client.AccountName, input.PolicyID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetPolicy request: {{err}}", err)
	}

	var result *Policy
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetPolicy response: {{err}}", err)
	}

	return result, nil
}

type DeletePolicyInput struct {
	PolicyID string
}

func (c *PoliciesClient) Delete(ctx context.Context, input *DeletePolicyInput) error {
	path := fmt.Sprintf("/%s/policies/%s", c.client.AccountName, input.PolicyID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeletePolicy request: {{err}}", err)
	}

	return nil
}

// UpdatePolicyInput represents the options that can be specified
// when updating a policy. Anything but ID can be modified.
type UpdatePolicyInput struct {
	PolicyID    string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Rules       []string `json:"rules,omitempty"`
	Description string   `json:"description,omitempty"`
}

func (c *PoliciesClient) Update(ctx context.Context, input *UpdatePolicyInput) (*Policy, error) {
	path := fmt.Sprintf("/%s/policies/%s", c.client.AccountName, input.PolicyID)
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
		return nil, errwrap.Wrapf("Error executing UpdatePolicy request: {{err}}", err)
	}

	var result *Policy
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding UpdatePolicy response: {{err}}", err)
	}

	return result, nil
}

type CreatePolicyInput struct {
	Name        string   `json:"name"`
	Rules       []string `json:"rules"`
	Description string   `json:"description,omitempty"`
}

func (c *PoliciesClient) Create(ctx context.Context, input *CreatePolicyInput) (*Policy, error) {
	path := fmt.Sprintf("/%s/policies", c.client.AccountName)
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
		return nil, errwrap.Wrapf("Error executing CreatePolicy request: {{err}}", err)
	}

	var result *Policy
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding CreatePolicy response: {{err}}", err)
	}

	return result, nil
}
