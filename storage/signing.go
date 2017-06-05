package storage

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
)

// SignURLInput represents parameters to a SignURL operation.
type SignURLInput struct {
	ValidityPeriod time.Duration
	Method         string
	ObjectPath     string
}

// SignURLOutput contains the outputs of a SignURL operation. To simply
// access the signed URL, use the SignedURL method.
type SignURLOutput struct {
	host       string
	objectPath string
	Method     string
	Algorithm  string
	Signature  string
	Expires    string
	KeyID      string
}

// SignedURL returns a signed URL for the given scheme. Valid schemes are
// `http` and `https`.
func (output *SignURLOutput) SignedURL(scheme string) string {
	query := &url.Values{}
	query.Set("algorithm", output.Algorithm)
	query.Set("expires", output.Expires)
	query.Set("keyId", output.KeyID)
	query.Set("signature", output.Signature)

	sUrl := url.URL{}
	sUrl.Scheme = scheme
	sUrl.Host = output.host
	sUrl.Path = output.objectPath
	sUrl.RawQuery = query.Encode()

	return sUrl.String()
}

// SignURL creates a time-expiring URL that can be shared with others.
// This is useful to generate HTML links, for example.
func (s *Storage) SignURL(input *SignURLInput) (*SignURLOutput, error) {
	hostUrl, err := url.Parse(s.client.Endpoint)
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing endpoint URL: {{err}}", err)
	}

	output := &SignURLOutput{
		host:       hostUrl.Host,
		objectPath: fmt.Sprintf("%s/stor/%s", s.client.AccountName, input.ObjectPath),
		Method:     input.Method,
		Algorithm:  strings.ToUpper(s.client.Authorizers[0].DefaultAlgorithm()),
		Expires:    strconv.FormatInt(time.Now().Add(input.ValidityPeriod).Unix(), 10),
		KeyID:      fmt.Sprintf("/%s/keys/%s", s.client.AccountName, s.client.Authorizers[0].KeyFingerprint()),
	}

	toSign := bytes.Buffer{}
	toSign.WriteString(input.Method + "\n")
	toSign.WriteString(hostUrl.Host + "\n")
	toSign.WriteString(fmt.Sprintf("/%s/stor/%s\n", s.client.AccountName, input.ObjectPath))

	query := &url.Values{}
	query.Set("algorithm", output.Algorithm)
	query.Set("expires", output.Expires)
	query.Set("keyId", output.KeyID)
	toSign.WriteString(query.Encode())

	signature, _, err := s.client.Authorizers[0].SignRaw(toSign.String())
	if err != nil {
		return nil, errwrap.Wrapf("Error signing string: {{err}}", err)
	}

	output.Signature = signature
	return output, nil
}
