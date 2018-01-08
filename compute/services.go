package compute

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"sort"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type ServicesClient struct {
	client *client.Client
}

type Service struct {
	Name     string
	Endpoint string
}

type ListServicesInput struct{}

func (c *ServicesClient) List(ctx context.Context, _ *ListServicesInput) ([]*Service, error) {
	fullPath := path.Join("/", c.client.AccountName, "services")
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to list services")
	}

	var intermediate map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&intermediate); err != nil {
		return nil, errors.Wrap(err, "unable to decode list services response")
	}

	keys := make([]string, len(intermediate))
	i := 0
	for k := range intermediate {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	result := make([]*Service, len(intermediate))
	i = 0
	for _, key := range keys {
		result[i] = &Service{
			Name:     key,
			Endpoint: intermediate[key],
		}
		i++
	}

	return result, nil
}
