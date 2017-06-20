package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type ImagesClient struct {
	client *client.Client
}

type ImageFile struct {
	Compression string `json:"compression"`
	SHA1        string `json:"sha1"`
	Size        int64  `json:"size"`
}

type Image struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	OS           string                 `json:"os"`
	Description  string                 `json:"description"`
	Version      string                 `json:"version"`
	Type         string                 `json:"type"`
	Requirements map[string]interface{} `json:"requirements"`
	Homepage     string                 `json:"homepage"`
	Files        []*ImageFile           `json:"files"`
	PublishedAt  time.Time              `json:"published_at"`
	Owner        string                 `json:"owner"`
	Public       bool                   `json:"public"`
	State        string                 `json:"state"`
	Tags         map[string]string      `json:"tags"`
	EULA         string                 `json:"eula"`
	ACL          []string               `json:"acl"`
	Error        TritonError            `json:"error"`
}

type ListImagesInput struct{}

func (c *ImagesClient) ListImages(ctx context.Context, _ *ListImagesInput) ([]*Image, error) {
	path := fmt.Sprintf("/%s/images", c.client.AccountName)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListImages request: {{err}}", err)
	}

	var result []*Image
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListImages response: {{err}}", err)
	}

	return result, nil
}

type GetImageInput struct {
	ImageID string
}

func (c *ImagesClient) GetImage(ctx context.Context, input *GetImageInput) (*Image, error) {
	path := fmt.Sprintf("/%s/images/%s", c.client.AccountName, input.ImageID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetImage request: {{err}}", err)
	}

	var result *Image
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetImage response: {{err}}", err)
	}

	return result, nil
}

type DeleteImageInput struct {
	ImageID string
}

func (c *ImagesClient) DeleteImage(ctx context.Context, input *DeleteImageInput) error {
	path := fmt.Sprintf("/%s/images/%s", c.client.AccountName, input.ImageID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteKey request: {{err}}", err)
	}

	return nil
}

type ExportImageInput struct {
	ImageID   string
	MantaPath string
}

type MantaLocation struct {
	MantaURL     string `json:"manta_url"`
	ImagePath    string `json:"image_path"`
	ManifestPath string `json:"manifest_path"`
}

func (c *ImagesClient) ExportImage(ctx context.Context, input *ExportImageInput) (*MantaLocation, error) {
	path := fmt.Sprintf("/%s/images/%s", c.client.AccountName, input.ImageID)
	query := &url.Values{}
	query.Set("action", "export")
	query.Set("manta_path", input.MantaPath)

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
		Query:  query,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetImage request: {{err}}", err)
	}

	var result *MantaLocation
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetImage response: {{err}}", err)
	}

	return result, nil
}

type CreateImageFromMachineInput struct {
	MachineID   string            `json:"machine"`
	Name        string            `json:"name"`
	Version     string            `json:"version,omitempty"`
	Description string            `json:"description,omitempty"`
	HomePage    string            `json:"homepage,omitempty"`
	EULA        string            `json:"eula,omitempty"`
	ACL         []string          `json:"acl,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

func (c *ImagesClient) CreateImageFromMachine(ctx context.Context, input *CreateImageFromMachineInput) (*Image, error) {
	path := fmt.Sprintf("/%s/images", c.client.AccountName)
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
		return nil, errwrap.Wrapf("Error executing CreateImageFromMachine request: {{err}}", err)
	}

	var result *Image
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding CreateImageFromMachine response: {{err}}", err)
	}

	return result, nil
}

type UpdateImageInput struct {
	ImageID     string            `json:"-"`
	Name        string            `json:"name"`
	Version     string            `json:"version,omitempty"`
	Description string            `json:"description,omitempty"`
	HomePage    string            `json:"homepage,omitempty"`
	EULA        string            `json:"eula,omitempty"`
	ACL         []string          `json:"acl,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

func (c *ImagesClient) UpdateImage(ctx context.Context, input *UpdateImageInput) (*Image, error) {
	path := fmt.Sprintf("/%s/images/%s", c.client.AccountName, input.ImageID)
	query := &url.Values{}
	query.Set("action", "update")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Query:  query,
		Body:   input,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing UpdateImage request: {{err}}", err)
	}

	var result *Image
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding UpdateImage response: {{err}}", err)
	}

	return result, nil
}
