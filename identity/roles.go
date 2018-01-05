package identity

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type RolesClient struct {
	client *client.Client
}

type Role struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Policies       []string `json:"policies"`
	Members        []string `json:"members"`
	DefaultMembers []string `json:"default_members"`
}

type ListRolesInput struct{}

func (c *RolesClient) List(ctx context.Context, _ *ListRolesInput) ([]*Role, error) {
	fullPath := path.Join("/", c.client.AccountName, "roles")

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListRoles request: {{err}}", err)
	}

	var result []*Role
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListRoles response: {{err}}", err)
	}

	return result, nil
}

type GetRoleInput struct {
	RoleID string
}

func (c *RolesClient) Get(ctx context.Context, input *GetRoleInput) (*Role, error) {
	fullPath := path.Join("/", c.client.AccountName, "roles", input.RoleID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetRole request: {{err}}", err)
	}

	var result *Role
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetRole response: {{err}}", err)
	}

	return result, nil
}

// CreateRoleInput represents the options that can be specified
// when creating a new role.
type CreateRoleInput struct {
	// Name of the role. Required.
	Name string `json:"name"`

	// This account's policies to be given to this role. Optional.
	Policies []string `json:"policies,omitempty"`

	// This account's user logins to be added to this role. Optional.
	Members []string `json:"members,omitempty"`

	// This account's user logins to be added to this role and have
	// it enabled by default. Optional.
	DefaultMembers []string `json:"default_members,omitempty"`
}

func (c *RolesClient) Create(ctx context.Context, input *CreateRoleInput) (*Role, error) {
	fullPath := path.Join("/", c.client.AccountName, "roles")
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing CreateRole request: {{err}}", err)
	}

	var result *Role
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding CreateRole response: {{err}}", err)
	}

	return result, nil
}

// UpdateRoleInput represents the options that can be specified
// when updating a role. Anything but ID can be modified.
type UpdateRoleInput struct {
	// ID of the role to modify. Required.
	RoleID string `json:"id"`

	// Name of the role. Required.
	Name string `json:"name"`

	// This account's policies to be given to this role. Optional.
	Policies []string `json:"policies,omitempty"`

	// This account's user logins to be added to this role. Optional.
	Members []string `json:"members,omitempty"`

	// This account's user logins to be added to this role and have
	// it enabled by default. Optional.
	DefaultMembers []string `json:"default_members,omitempty"`
}

func (c *RolesClient) Update(ctx context.Context, input *UpdateRoleInput) (*Role, error) {
	fullPath := path.Join("/", c.client.AccountName, "roles", input.RoleID)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing UpdateRole request: {{err}}", err)
	}

	var result *Role
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding UpdateRole response: {{err}}", err)
	}

	return result, nil
}

type DeleteRoleInput struct {
	RoleID string
}

func (c *RolesClient) Delete(ctx context.Context, input *DeleteRoleInput) error {
	fullPath := path.Join("/", c.client.AccountName, "roles", input.RoleID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteRole request: {{err}}", err)
	}

	return nil
}
