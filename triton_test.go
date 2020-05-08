//
// Copyright 2020 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package triton_test

import (
	"os"
	"testing"

	triton "github.com/joyent/triton-go/v2"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name    string
		varname string
		input   string
		value   string
	}{
		{"Triton", "TRITON_NAME", "NAME", "good"},
		{"SDC", "SDC_NAME", "NAME", "good"},
		{"unrelated", "BAD_NAME", "NAME", ""},
		{"missing", "", "NAME", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.varname, test.value)
			defer os.Unsetenv(test.varname)

			if val := triton.GetEnv(test.input); val != test.value {
				t.Errorf("expected %s env var to be '%s': got '%s'",
					test.varname, test.value, val)
			}
		})
	}
}
