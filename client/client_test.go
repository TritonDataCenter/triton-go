//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package client

import (
	"os"
	"strings"
	"testing"

	auth "github.com/joyent/triton-go/authentication"
)

const BadURL = "**ftp://man($$"

// DON'T USE THIS OUTSIDE TESTING ~ This key was only created to use for
// internal unit testing. It should never be used for acceptance testing either.
//
// This is just a randomly generated key pair.
var DummyAuth = struct {
	Fingerprint string
	PrivateKey  []byte
	PublicKey   []byte
	Signer      auth.Signer
}{
	"9f:d6:65:fc:d6:60:dc:d0:4e:db:2d:75:f7:92:8c:31",
	[]byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`),
	[]byte(`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC6L2U2MIlqEd4VIULodeL8I1fnzwL8vaAlTHrluc2mFUL3yFbXDOhG89cfHlpadvuURuk2E/oHRSgIENNo+C32Ed8P+g1D7odPVCkMgZ9mZcj3X+1yMZVIFvE3tZpi38ew+saNtqXC9A8KE3newxWYWQxMQNHkbfM4E4Aas6KaUJPzTAjZLleDSyl3yOMozkjWzANHTG24NacdNHU6GKf6kRAgTqCUtN9iruhDHqeYxdcE2aT5Jo8PUUDKOQX/zWKbYQdcHKA0TCtpWsNwDEraYqVTZcwmKiBAo1haWbdSRs0HpY8t7Xb4prTauYhWxKoilz9kIGtoCUJpJaVgKRT68naEf/7KnLltc2VtpOJZkr6+4js8LztPnTIs15dW41kDKfMOkEK+yOKIesH9kVF9V8Xo++JXwJX6bsZPN8DAmFV5JuWtHqe2zXorQdMPjy/nsNTRvS3lU68LtxY9KGz5p91IvdnxF4vkRYyBTXB++3JxEp6GINUaPcyzZx/zdEgA2vfzsb7UtxNBg/skm4HtK/LmP5Hr4eMjRyC7gQxQ91Z8I0YbKuXHAk20+ZRJ9qwfXFQqmqvt0L7Sw3j+B8nDvlZMBisBOQ/N7Ghg1oywYN+G6iGRhrdOfarWDeSqGxf9HQV3l3x3mpzXfupxELGJljP1/dtWRdBuFiY711rcqw== test-dummy-20171002140848`),
	nil,
}

func TestNew(t *testing.T) {
	mantaURL := "https://us-east.manta.joyent.com"
	tsgEnv := "http://tsg.test.org"
	jpcTritonURL := "https://us-east-1.api.joyent.com"
	spcTritonURL := "https://us-east-1.api.samsungcloud.io"
	jpcServiceURL := "https://tsg.us-east-1.svc.joyent.zone"
	spcServiceURL := "https://tsg.us-east-1.svc.samsungcloud.zone"
	privateInstallUrl := "https://myinstall.mycompany.com"

	accountName := "test.user"
	signer, _ := auth.NewTestSigner()

	tests := []struct {
		name        string
		tritonURL   string
		mantaURL    string
		tsgEnv      string
		servicesURL string
		accountName string
		signer      auth.Signer
		err         interface{}
	}{
		{"default", jpcTritonURL, mantaURL, "", jpcServiceURL, accountName, signer, nil},
		{"in samsung", spcTritonURL, mantaURL, "", spcServiceURL, accountName, signer, nil},
		{"env TSG", jpcTritonURL, mantaURL, tsgEnv, tsgEnv, accountName, signer, nil},
		{"missing url", "", "", "", "", accountName, signer, ErrMissingURL},
		{"bad tritonURL", BadURL, mantaURL, "", "", accountName, signer, InvalidTritonURL},
		{"bad mantaURL", jpcTritonURL, BadURL, "", jpcServiceURL, accountName, signer, InvalidMantaURL},
		{"bad TSG", jpcTritonURL, mantaURL, BadURL, "", accountName, signer, InvalidServicesURL},
		{"missing accountName", jpcTritonURL, mantaURL, "", jpcServiceURL, "", signer, ErrAccountName},
		{"missing signer", jpcTritonURL, mantaURL, "", jpcServiceURL, accountName, nil, ErrDefaultAuth},
		{"private install", privateInstallUrl, mantaURL, "", "", accountName, signer, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Unsetenv("TRITON_KEY_ID")
			os.Unsetenv("SDC_KEY_ID")
			os.Unsetenv("MANTA_KEY_ID")
			os.Unsetenv("SSH_AUTH_SOCK")
			os.Unsetenv("TRITON_TSG_URL")

			if test.tsgEnv != "" {
				os.Setenv("TRITON_TSG_URL", test.tsgEnv)
			}

			c, err := New(
				test.tritonURL,
				test.mantaURL,
				test.accountName,
				test.signer,
			)

			// test generation of TSG URL for all non-error cases
			if err == nil {
				if c.ServicesURL.String() != test.servicesURL {
					t.Errorf("expected ServicesURL to be set to %q: got %q (%s)",
						test.servicesURL, c.ServicesURL.String(), test.name)
					return
				}
			}

			if test.err != nil {
				if err == nil {
					t.Error("expected error not to be nil")
					return
				}

				switch test.err.(type) {
				case error:
					testErr := test.err.(error)
					if err.Error() != testErr.Error() {
						t.Errorf("expected error: received %v", err)
					}
				case string:
					testErr := test.err.(string)
					if !strings.Contains(err.Error(), testErr) {
						t.Errorf("expected error: received %v", err)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("expected error to be nil: received %v", err)
			}
		})
	}

	t.Run("default SSH agent auth", func(t *testing.T) {
		os.Unsetenv("SSH_AUTH_SOCK")
		err := os.Setenv("TRITON_KEY_ID", DummyAuth.Fingerprint)
		defer os.Unsetenv("TRITON_KEY_ID")
		if err != nil {
			t.Errorf("expected error to not be nil: received %v", err)
		}

		_, err = New(
			jpcTritonURL,
			mantaURL,
			accountName,
			nil,
		)
		if err == nil {
			t.Error("expected error to not be nil")
		}
		if !strings.Contains(err.Error(), "unable to initialize NewSSHAgentSigner") {
			t.Errorf("expected error to be from NewSSHAgentSigner: received '%v'", err)
		}
	})
}
