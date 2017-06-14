package storage

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/errwrap"
)

// PutSnapLinkInput represents parameters to a PutSnapLink operation.
type PutSnapLinkInput struct {
	LinkPath   string
	SourcePath string
}

// PutSnapLink creates a SnapLink to an object.
func (s *Storage) PutSnapLink(input *PutSnapLinkInput) error {
	path := fmt.Sprintf("/%s/stor/%s", s.Client.AccountName, input.LinkPath)
	headers := &http.Header{}
	headers.Set("Content-Type", "application/json; type=link")
	headers.Set("Location", input.SourcePath)

	respBody, _, err := s.executeRequest(http.MethodPut, path, nil, headers, nil)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing PutSnapLink request: {{err}}", err)
	}

	return nil
}
