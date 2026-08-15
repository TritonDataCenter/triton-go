//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
// Copyright 2026 Edgecast Cloud LLC.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package authentication

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/md5"
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// formatPublicKeyFingerprint produces the MD5 fingerprint of the given SSH
// public key. If display is true, the fingerprint is formatted with colons
// between each byte, as per the output of OpenSSL.
func formatPublicKeyFingerprint(privateKey interface{}, display bool) (string, error) {
	var key ssh.PublicKey
	switch k := privateKey.(type) {
	case ssh.PublicKey:
		key = k
	case *rsa.PrivateKey:
		p, err := ssh.NewPublicKey(k.Public())
		if err != nil {
			return "", errors.Wrap(err, "unable to parse SSH key from private key")
		}
		key = p
	case *ecdsa.PrivateKey:
		p, err := ssh.NewPublicKey(k.Public())
		if err != nil {
			return "", errors.Wrap(err, "unable to parse SSH key from private key")
		}
		key = p
	case ed25519.PrivateKey:
		p, err := ssh.NewPublicKey(k.Public())
		if err != nil {
			return "", errors.Wrap(err, "unable to parse SSH key from private key")
		}
		key = p
	default:
		return "", fmt.Errorf("unable to parse SSH key from private key")
	}
	publicKeyFingerprint := md5.New()
	publicKeyFingerprint.Write(key.Marshal())
	publicKeyFingerprintString := fmt.Sprintf("%x", publicKeyFingerprint.Sum(nil))

	if !display {
		return publicKeyFingerprintString, nil
	}

	formatted := ""
	for i := 0; i < len(publicKeyFingerprintString); i = i + 2 {
		formatted = fmt.Sprintf("%s%s:", formatted, publicKeyFingerprintString[i:i+2])
	}

	return strings.TrimSuffix(formatted, ":"), nil
}
