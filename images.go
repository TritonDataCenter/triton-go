package triton

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/errwrap"
	"net/http"
	"time"
)

type ImagesClient struct {
	*Client
}

// Images returns a client used for accessing functions pertaining to
// Images functionality in the Triton API.
func (client *Client) Images() *ImagesClient {
	return &ImagesClient{client}
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

func (client *ImagesClient) ListImages(*ListImagesInput) ([]*Image, error) {
	respReader, err := client.executeRequest(http.MethodGet, "/my/images", nil)
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

func (client *ImagesClient) GetImage(input *GetImageInput) (*Image, error) {
	path := fmt.Sprintf("/%s/images/%s", client.accountName, input.ImageID)
	respReader, err := client.executeRequest(http.MethodGet, path, nil)
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
