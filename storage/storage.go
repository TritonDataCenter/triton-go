package storage

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type Storage struct {
	Client *client.Client
}

func (s *Storage) executeRequest(method, path string, query *url.Values, headers *http.Header, body interface{}) (io.ReadCloser, http.Header, error) {
	var requestBody io.ReadSeeker
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	endpoint, err := url.Parse(os.Getenv("MANTA_URL"))
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error parsing MANTA_URL: {{err}}", err)
	}
	endpoint.Path = path

	req, err := http.NewRequest(method, endpoint.String(), requestBody)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error constructing HTTP request: {{err}}", err)
	}

	if body != nil && (headers == nil || headers.Get("Content-Type") == "") {
		req.Header.Set("Content-Type", "application/json")
	}
	if headers != nil {
		for key, values := range *headers {
			for _, value := range values {
				req.Header.Set(key, value)
			}
		}
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	authHeader, err := s.Client.Authorizers[0].Sign(dateHeader)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error signing HTTP request: {{err}}", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "manta-go client API")

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	resp, err := s.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error executing HTTP request: {{err}}", err)
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, resp.Header, nil
	}

	mantaError := &MantaError{
		StatusCode: resp.StatusCode,
	}

	errorDecoder := json.NewDecoder(resp.Body)
	if err := errorDecoder.Decode(mantaError); err != nil {
		return nil, nil, errwrap.Wrapf("Error decoding error response: {{err}}", err)
	}
	return nil, nil, mantaError
}

func (s *Storage) executeRequestNoEncode(method, path string, query *url.Values, headers *http.Header, body io.ReadSeeker) (io.ReadCloser, http.Header, error) {
	endpoint, err := url.Parse(os.Getenv("MANTA_URL"))
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error parsing MANTA_URL: {{err}}", err)
	}
	endpoint.Path = path

	req, err := http.NewRequest(method, endpoint.String(), body)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error constructing HTTP request: {{err}}", err)
	}

	if headers != nil {
		for key, values := range *headers {
			for _, value := range values {
				req.Header.Set(key, value)
			}
		}
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	authHeader, err := s.Client.Authorizers[0].Sign(dateHeader)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error signing HTTP request: {{err}}", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "manta-go client API")

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	resp, err := s.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, errwrap.Wrapf("Error executing HTTP request: {{err}}", err)
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, resp.Header, nil
	}

	mantaError := &MantaError{
		StatusCode: resp.StatusCode,
	}

	errorDecoder := json.NewDecoder(resp.Body)
	if err := errorDecoder.Decode(mantaError); err != nil {
		return nil, nil, errwrap.Wrapf("Error decoding error response: {{err}}", err)
	}
	return nil, nil, mantaError
}
