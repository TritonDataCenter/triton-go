package compute

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"context"

	"github.com/hashicorp/errwrap"
)

type DataCentersClient struct {
	*Compute
}

type DataCenter struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ListDataCentersInput struct{}

func (c *DataCentersClient) ListDataCenters(ctx context.Context, _ *ListDataCentersInput) ([]*DataCenter, error) {
	path := fmt.Sprintf("/%s/datacenters", c.client.AccountName)
	respReader, err := c.executeRequest(ctx, http.MethodGet, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListDatacenters request: {{err}}", err)
	}

	var intermediate map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&intermediate); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListDatacenters response: {{err}}", err)
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

func (c *DataCentersClient) GetDataCenter(ctx context.Context, input *GetDataCenterInput) (*DataCenter, error) {
	path := fmt.Sprintf("/%s/datacenters/%s", c.client.AccountName, input.Name)
	resp, err := c.executeRequestRaw(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetDatacenter request: {{err}}", err)
	}

	if resp.StatusCode != http.StatusFound {
		return nil, fmt.Errorf("Error executing GetDatacenter request: expected status code 302, got %s",
			resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return nil, errors.New("Error decoding GetDatacenter response: no Location header")
	}

	return &DataCenter{
		Name: input.Name,
		URL:  location,
	}, nil
}
