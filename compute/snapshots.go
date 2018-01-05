package compute

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type SnapshotsClient struct {
	client *client.Client
}

type Snapshot struct {
	Name    string
	State   string
	Created time.Time
	Updated time.Time
}

type ListSnapshotsInput struct {
	MachineID string
}

func (c *SnapshotsClient) List(ctx context.Context, input *ListSnapshotsInput) ([]*Snapshot, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots")
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing List request: {{err}}", err)
	}

	var result []*Snapshot
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding List response: {{err}}", err)
	}

	return result, nil
}

type GetSnapshotInput struct {
	MachineID string
	Name      string
}

func (c *SnapshotsClient) Get(ctx context.Context, input *GetSnapshotInput) (*Snapshot, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing Get request: {{err}}", err)
	}

	var result *Snapshot
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding Get response: {{err}}", err)
	}

	return result, nil
}

type DeleteSnapshotInput struct {
	MachineID string
	Name      string
}

func (c *SnapshotsClient) Delete(ctx context.Context, input *DeleteSnapshotInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing Delete request: {{err}}", err)
	}

	return nil
}

type StartMachineFromSnapshotInput struct {
	MachineID string
	Name      string
}

func (c *SnapshotsClient) StartMachine(ctx context.Context, input *StartMachineFromSnapshotInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing StartMachine request: {{err}}", err)
	}

	return nil
}

type CreateSnapshotInput struct {
	MachineID string
	Name      string
}

func (c *SnapshotsClient) Create(ctx context.Context, input *CreateSnapshotInput) (*Snapshot, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots")

	data := make(map[string]interface{})
	data["name"] = input.Name

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   data,
	}

	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing Create request: {{err}}", err)
	}

	var result *Snapshot
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding Create response: {{err}}", err)
	}

	return result, nil
}
