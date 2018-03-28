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

func TestNew(t *testing.T) {
	tritonURL := "https://us-east-1.api.joyent.com"
	mantaURL := "https://us-east.manta.joyent.com"
	servicesURL := "https://tsg.us-east-1.svc.joyent.zone"
	tsgEnv := "http://tsg.test.org"
	accountName := "test.user"
	signer, _ := auth.NewTestSigner()

	tests := []struct {
		name        string
		tritonURL   string
		mantaURL    string
		tsgEnv      string
		accountName string
		signer      auth.Signer
		err         interface{}
	}{
		{"default", tritonURL, mantaURL, "", accountName, signer, nil},
		{"env TSG", tritonURL, mantaURL, tsgEnv, accountName, signer, nil},
		{"missing url", "", "", "", accountName, signer, ErrMissingURL},
		{"bad tritonURL", BadURL, mantaURL, "", accountName, signer, InvalidTritonURL},
		{"bad mantaURL", tritonURL, BadURL, "", accountName, signer, InvalidMantaURL},
		{"bad TSG", tritonURL, mantaURL, BadURL, accountName, signer, InvalidServicesURL},
		{"missing accountName", tritonURL, mantaURL, "", "", signer, ErrAccountName},
		{"missing signer", tritonURL, mantaURL, "", accountName, nil, ErrDefaultAuth},
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

			// NOTE: test the generation of our TSG URL for all non-error cases
			if err == nil {
				if test.tsgEnv == "" {
					if c.ServicesURL.String() != servicesURL {
						t.Errorf("expected ServicesURL to be set to %s: got %s",
							servicesURL, c.ServicesURL.String())
						return
					}
				} else {
					if c.ServicesURL.String() != tsgEnv {
						t.Errorf("expected ServicesURL to be set to %s: got %s",
							tsgEnv, c.ServicesURL.String())
						return
					}
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
		err := os.Setenv("TRITON_KEY_ID", auth.Dummy.Fingerprint)
		defer os.Unsetenv("TRITON_KEY_ID")
		if err != nil {
			t.Errorf("expected error to not be nil: received %v", err)
		}

		_, err = New(
			tritonURL,
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
