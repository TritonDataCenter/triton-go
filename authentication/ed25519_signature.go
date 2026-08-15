//
// Copyright 2026 Edgecast Cloud LLC.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package authentication

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

type ed25519Signature struct {
	signature []byte
}

func (s *ed25519Signature) SignatureType() string {
	return ED25519_SHA512
}

func (s *ed25519Signature) String() string {
	return base64.StdEncoding.EncodeToString(s.signature)
}

func newEd25519Signature(signatureBlob []byte) (*ed25519Signature, error) {
	if len(signatureBlob) != ed25519.SignatureSize {
		return nil, fmt.Errorf("invalid ed25519 signature length: got %d, want %d",
			len(signatureBlob), ed25519.SignatureSize)
	}
	return &ed25519Signature{
		signature: signatureBlob,
	}, nil
}
