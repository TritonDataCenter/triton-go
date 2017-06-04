package manta

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
)

// GetObjectInput represents parameters to a GetObject operation.
type GetObjectInput struct {
	ObjectPath string
}

// GetObjectOutput contains the outputs for a GetObject operation. It is your
// responsibility to ensure that the io.ReadCloser ObjectReader is closed.
type GetObjectOutput struct {
	ContentLength uint64
	ContentType   string
	LastModified  time.Time
	ContentMD5    string
	ETag          string
	Metadata      map[string]string
	ObjectReader  io.ReadCloser
}

// GetObject retrieves an object from the Manta service. If error is nil (i.e.
// the call returns successfully), it is your responsibility to close the io.ReadCloser
// named ObjectReader in the operation output.
func (c *Client) GetObject(input *GetObjectInput) (*GetObjectOutput, error) {
	path := fmt.Sprintf("/%s/stor/%s", c.accountName, input.ObjectPath)

	respBody, respHeaders, err := c.executeRequest(http.MethodGet, path, nil, nil, nil)
	if err != nil {
		respBody.Close()
		return nil, errwrap.Wrapf("Error executing GetDirectory request: {{err}}", err)
	}

	response := &GetObjectOutput{
		ContentType:  respHeaders.Get("Content-Type"),
		ContentMD5:   respHeaders.Get("Content-MD5"),
		ETag:         respHeaders.Get("Etag"),
		ObjectReader: respBody,
	}

	lastModified, err := time.Parse(time.RFC1123, respHeaders.Get("Last-Modified"))
	if err == nil {
		response.LastModified = lastModified
	}

	contentLength, err := strconv.ParseUint(respHeaders.Get("Content-Length"), 10, 64)
	if err == nil {
		response.ContentLength = contentLength
	}

	metadata := map[string]string{}
	for key, values := range respHeaders {
		if strings.HasPrefix(key, "m-") {
			metadata[key] = strings.Join(values, ", ")
		}
	}
	response.Metadata = metadata

	return response, nil
}

// DeleteObjectInput represents parameters to a DeleteObject operation.
type DeleteObjectInput struct {
	ObjectPath string
}

// DeleteObject deletes an object.
func (c *Client) DeleteObject(input *DeleteObjectInput) error {
	path := fmt.Sprintf("/%s/stor/%s", c.accountName, input.ObjectPath)

	respBody, _, err := c.executeRequest(http.MethodDelete, path, nil, nil, nil)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteObject request: {{err}}", err)
	}

	return nil
}

// PutObjectMetadataInput represents parameters to a PutObjectMetadata operation.
type PutObjectMetadataInput struct {
	ObjectPath  string
	ContentType string
	Metadata    map[string]string
}

// PutObjectMetadata allows you to overwrite the HTTP headers for an already
// existing object, without changing the data. Note this is an idempotent "replace"
// operation, so you must specify the complete set of HTTP headers you want
// stored on each request.
//
// You cannot change "critical" headers:
// 	- Content-Length
//	- Content-MD5
//	- Durability-Level
func (c *Client) PutObjectMetadata(input *PutObjectMetadataInput) error {
	path := fmt.Sprintf("/%s/stor/%s", c.accountName, input.ObjectPath)
	query := &url.Values{}
	query.Set("metadata", "true")

	headers := &http.Header{}
	headers.Set("Content-Type", input.ContentType)
	for key, value := range input.Metadata {
		headers.Set(key, value)
	}

	respBody, _, err := c.executeRequest(http.MethodPut, path, query, headers, nil)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing PutObjectMetadata request: {{err}}", err)
	}

	return nil
}

// PutObjectInput represents parameters to a PutObject operation.
type PutObjectInput struct {
	ObjectPath       string
	DurabilityLevel  uint64
	ContentType      string
	ContentMD5       string
	IfMatch          string
	IfModifiedSince  *time.Time
	ContentLength    uint64
	MaxContentLength uint64
	ObjectReader     io.ReadSeeker
}

func (c *Client) PutObject(input *PutObjectInput) error {
	path := fmt.Sprintf("/%s/stor/%s", c.accountName, input.ObjectPath)

	if input.MaxContentLength != 0 && input.ContentLength != 0 {
		return errors.New("ContentLength and MaxContentLength may not both be set to non-zero values.")
	}

	headers := &http.Header{}
	if input.DurabilityLevel != 0 {
		headers.Set("Durability-Level", strconv.FormatUint(input.DurabilityLevel, 10))
	}
	if input.ContentType != "" {
		headers.Set("Content-Type", input.ContentType)
	}
	if input.ContentMD5 != "" {
		headers.Set("Content-MD$", input.ContentMD5)
	}
	if input.IfMatch != "" {
		headers.Set("If-Match", input.IfMatch)
	}
	if input.IfModifiedSince != nil {
		headers.Set("If-Modified-Since", input.IfModifiedSince.Format(time.RFC1123))
	}
	if input.ContentLength != 0 {
		headers.Set("Content-Length", strconv.FormatUint(input.ContentLength, 10))
	}
	if input.MaxContentLength != 0 {
		headers.Set("Max-Content-Length", strconv.FormatUint(input.MaxContentLength, 10))
	}

	respBody, _, err := c.executeRequestNoEncode(http.MethodPut, path, nil, headers, input.ObjectReader)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing PutObjectMetadata request: {{err}}", err)
	}

	return nil
}
