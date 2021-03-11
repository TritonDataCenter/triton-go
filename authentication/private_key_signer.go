//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package authentication

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

const (
	RSA_SHA512     = "rsa-sha512"
	ECDSA_SHA512   = "ecdsa-sha512"
	DSA_SHA512     = "dsa-sha512"
	ED25519_SHA512 = "ed25519-sha512"
)

type PrivateKeySigner struct {
	formattedKeyFingerprint string
	keyFingerprint          string
	algorithm               string
	accountName             string
	userName                string
	hashFunc                crypto.Hash

	privateKey interface{}
}

type PrivateKeySignerInput struct {
	KeyID              string
	PrivateKeyMaterial []byte
	AccountName        string
	Username           string
}

func NewPrivateKeySigner(input PrivateKeySignerInput) (*PrivateKeySigner, error) {
	keyFingerprintMD5 := strings.Replace(input.KeyID, ":", "", -1)

	key, err := ssh.ParseRawPrivateKey(input.PrivateKeyMaterial)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse private key")
	}

	matchKeyFingerprint, err := formatPublicKeyFingerprint(key, false)
	if err != nil {
		return nil, errors.Wrap(err, "unable to format match public key")
	}
	displayKeyFingerprint, err := formatPublicKeyFingerprint(key, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to format display public key")
	}
	if matchKeyFingerprint != keyFingerprintMD5 {
		return nil, errors.New("Private key file does not match public key fingerprint")
	}

	signer := &PrivateKeySigner{
		formattedKeyFingerprint: displayKeyFingerprint,
		keyFingerprint:          input.KeyID,
		accountName:             input.AccountName,

		hashFunc:   crypto.SHA512,
		privateKey: key,
	}

	if input.Username != "" {
		signer.userName = input.Username
	}

	_, algorithm, err := signer.SignRaw("HelloWorld")
	if err != nil {
		return nil, fmt.Errorf("Cannot sign using ssh agent: %s", err)
	}
	signer.algorithm = algorithm

	return signer, nil
}

func (s *PrivateKeySigner) Sign(dateHeader string, isManta bool) (string, error) {
	const headerName = "date"

	hash := s.hashFunc.New()
	hash.Write([]byte(fmt.Sprintf("%s: %s", headerName, dateHeader)))
	digest := hash.Sum(nil)

	var algoName string
	var signedBase64 string
	switch s.privateKey.(type) {
	case *rsa.PrivateKey:
		algoName = RSA_SHA512
		signed, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey.(*rsa.PrivateKey), s.hashFunc, digest)
		if err != nil {
			return "", errors.Wrap(err, "unable to sign date header")
		}
		signedBase64 = base64.StdEncoding.EncodeToString(signed)
	case *ecdsa.PrivateKey:
		algoName = ECDSA_SHA512
		r, s, err := ecdsa.Sign(rand.Reader, s.privateKey.(*ecdsa.PrivateKey), digest)
		if err != nil {
			return "", errors.Wrap(err, "unable to sign date header")
		}
		signature := ECDSASignature{R: r, S: s}
		signed, err := asn1.Marshal(signature)
		signedBase64 = base64.StdEncoding.EncodeToString(signed)
	}

	key := &KeyID{
		UserName:    s.userName,
		AccountName: s.accountName,
		Fingerprint: s.formattedKeyFingerprint,
		IsManta:     isManta,
	}

	return fmt.Sprintf(authorizationHeaderFormat, key.generate(), algoName, headerName, signedBase64), nil
}

func (s *PrivateKeySigner) SignRaw(toSign string) (string, string, error) {
	hash := s.hashFunc.New()
	hash.Write([]byte(toSign))
	digest := hash.Sum(nil)

	var algoName string
	var signedBase64 string
	switch s.privateKey.(type) {
	case *rsa.PrivateKey:
		algoName = RSA_SHA512
		signed, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey.(*rsa.PrivateKey), s.hashFunc, digest)
		if err != nil {
			return "", "", errors.Wrap(err, "unable to sign date header")
		}
		signedBase64 = base64.StdEncoding.EncodeToString(signed)
	case *ecdsa.PrivateKey:
		algoName = ECDSA_SHA512
		r, s, err := ecdsa.Sign(rand.Reader, s.privateKey.(*ecdsa.PrivateKey), digest)
		if err != nil {
			return "", "", errors.Wrap(err, "unable to sign date header")
		}
		signature := ECDSASignature{R: r, S: s}
		signed, err := asn1.Marshal(signature)
		signedBase64 = base64.StdEncoding.EncodeToString(signed)

	}

	return signedBase64, algoName, nil
}

type ECDSASignature struct {
	R, S *big.Int
}

func (s *PrivateKeySigner) KeyFingerprint() string {
	return s.formattedKeyFingerprint
}

func (s *PrivateKeySigner) DefaultAlgorithm() string {
	return s.algorithm
}
