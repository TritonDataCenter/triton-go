package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "unable to list policies")
	}

	var result []*Policy
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode list policies response")
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
		return nil, errors.Wrap(err, "unable to get policy")
	}

	var result *Policy
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode get policy response")
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
		return errors.Wrap(err, "unable to delete policy")
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
		return nil, errors.Wrap(err, "unable to update policy")
	}

	var result *Policy
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode update policy response")
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
		return nil, errors.Wrap(err, "unable to create policy")
	}

	var result *Policy
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode create policy response")
	}

	return result, nil
}
