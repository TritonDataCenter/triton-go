package storage

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type SnapLinksClient struct {
	client *client.Client
}

// PutSnapLinkInput represents parameters to a PutSnapLink operation.
type PutSnapLinkInput struct {
	LinkPath   string
	SourcePath string
}

// PutSnapLink creates a SnapLink to an object.
func (s *SnapLinksClient) Put(input *PutSnapLinkInput) error {
	path := fmt.Sprintf("/%s%s", s.client.AccountName, input.LinkPath)
	headers := &http.Header{}
	headers.Set("Content-Type", "application/json; type=link")
	headers.Set("Location", input.SourcePath)

	reqInput := client.RequestInput{
		Method:  http.MethodPut,
		Path:    path,
		Headers: headers,
	}
	respBody, _, err := s.client.ExecuteRequestStorage(reqInput)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing PutSnapLink request: {{err}}", err)
	}

	return nil
}
