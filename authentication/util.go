//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package authentication

import (
	"crypto/ecdsa"
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
	switch privateKey.(type) {
	case *rsa.PrivateKey:
		p, err := ssh.NewPublicKey(privateKey.(*rsa.PrivateKey).Public())
		if err != nil {
			return "", errors.Wrap(err, "unable to parse SSH key from private key")
		}
		key = p
	case *ecdsa.PrivateKey:
		p, err := ssh.NewPublicKey(privateKey.(*ecdsa.PrivateKey).Public())
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
