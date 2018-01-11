package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type DataCentersClient struct {
	client *client.Client
}

type DataCenter struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ListDataCentersInput struct{}

func (c *DataCentersClient) List(ctx context.Context, _ *ListDataCentersInput) ([]*DataCenter, error) {
	fullPath := path.Join("/", c.client.AccountName, "datacenters")

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to list datacenters")
	}

	var intermediate map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&intermediate); err != nil {
		return nil, errors.Wrap(err, "unable to decode list datacenters response")
	}

	keys := make([]string, len(intermediate))
	i := 0
	for k := range intermediate {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	result := make([]*DataCenter, len(intermediate))
	i = 0
	for _, key := range keys {
		result[i] = &DataCenter{
			Name: key,
			URL:  intermediate[key],
		}
		i++
	}

	return result, nil
}

type GetDataCenterInput struct {
	Name string
}

func (c *DataCentersClient) Get(ctx context.Context, input *GetDataCenterInput) (*DataCenter, error) {
	fullPath := path.Join("/", c.client.AccountName, "datacenters", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	resp, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get datacenter")
	}

	if resp.StatusCode != http.StatusFound {
		return nil, fmt.Errorf("Error executing Get request: expected status code 302, got %d",
			resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return nil, errors.New("Error decoding Get response: no Location header")
	}

	return &DataCenter{
		Name: input.Name,
		URL:  location,
	}, nil
}
