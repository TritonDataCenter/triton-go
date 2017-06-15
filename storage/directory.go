package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/errwrap"
)

// DirectoryEntry represents an object or directory in Manta.
type DirectoryEntry struct {
	ETag         string    `json:"etag"`
	ModifiedTime time.Time `json:"mtime"`
	Name         string    `json:"name"`
	Size         uint64    `json:"size"`
	Type         string    `json:"type"`
}

// ListDirectoryInput represents parameters to a ListDirectory operation.
type ListDirectoryInput struct {
	DirectoryName string
	Limit         uint64
	Marker        string
}

// ListDirectoryOutput contains the outputs of a ListDirectory operation.
type ListDirectoryOutput struct {
	Entries       []*DirectoryEntry
	ResultSetSize uint64
}

// ListDirectory lists the contents of a directory.
func (s *Storage) ListDirectory(input *ListDirectoryInput) (*ListDirectoryOutput, error) {
	path := fmt.Sprintf("/%s%s", s.Client.AccountName, input.DirectoryName)
	query := &url.Values{}
	if input.Limit != 0 {
		query.Set("limit", strconv.FormatUint(input.Limit, 10))
	}
	if input.Marker != "" {
		query.Set("manta_path", input.Marker)
	}

	respBody, respHeader, err := s.executeRequest(http.MethodGet, path, query, nil, nil)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListDirectory request: {{err}}", err)
	}

	var results []*DirectoryEntry
	for {
		current := &DirectoryEntry{}
		decoder := json.NewDecoder(respBody)
		if err = decoder.Decode(&current); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errwrap.Wrapf("Error decoding ListDirectory response: {{err}}", err)
		}
		results = append(results, current)
	}

	output := &ListDirectoryOutput{
		Entries: results,
	}

	resultSetSize, err := strconv.ParseUint(respHeader.Get("Result-Set-Size"), 10, 64)
	if err == nil {
		output.ResultSetSize = resultSetSize
	}

	return output, nil
}

// PutDirectoryInput represents parameters to a PutDirectory operation.
type PutDirectoryInput struct {
	DirectoryName string
}

// PutDirectory in the Joyent Manta Storage Service is an idempotent create-or-update
// operation. Your private namespace starts at /:login/stor, and you can create any
// nested set of directories or objects underneath that.
func (s *Storage) PutDirectory(input *PutDirectoryInput) error {
	path := fmt.Sprintf("/%s/stor/%s", s.Client.AccountName, input.DirectoryName)
	headers := &http.Header{}
	headers.Set("Content-Type", "application/json; type=directory")

	respBody, _, err := s.executeRequest(http.MethodPut, path, nil, headers, nil)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing PutDirectory request: {{err}}", err)
	}

	return nil
}

// DeleteDirectoryInput represents parameters to a DeleteDirectory operation.
type DeleteDirectoryInput struct {
	DirectoryName string
}

// DeleteDirectory deletes a directory. The directory must be empty.
func (s *Storage) DeleteDirectory(input *DeleteDirectoryInput) error {
	path := fmt.Sprintf("/%s/stor/%s", s.Client.AccountName, input.DirectoryName)

	respBody, _, err := s.executeRequest(http.MethodDelete, path, nil, nil, nil)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteDirectory request: {{err}}", err)
	}

	return nil
}
