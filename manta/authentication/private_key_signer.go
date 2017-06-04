package authentication

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"golang.org/x/crypto/ssh"
)

type PrivateKeySigner struct {
	formattedKeyFingerprint string
	keyFingerprint          string
	accountName             string
	hashFunc                crypto.Hash

	privateKey *rsa.PrivateKey
}

func NewPrivateKeySigner(keyFingerprint string, privateKeyMaterial []byte, accountName string) (*PrivateKeySigner, error) {
	keyFingerprintMD5 := strings.Replace(keyFingerprint, ":", "", -1)

	block, _ := pem.Decode(privateKeyMaterial)
	if block == nil {
		return nil, errors.New("Error PEM-decoding private key material: nil block received")
	}

	rsakey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing private key: %s", err)
	}

	sshPublicKey, err := ssh.NewPublicKey(rsakey.Public())
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing SSH key from private key: %s", err)
	}

	matchKeyFingerprint := formatPublicKeyFingerprint(sshPublicKey, false)
	displayKeyFingerprint := formatPublicKeyFingerprint(sshPublicKey, true)
	if matchKeyFingerprint != keyFingerprintMD5 {
		return nil, errors.New("Private key file does not match public key fingerprint")
	}

	return &PrivateKeySigner{
		formattedKeyFingerprint: displayKeyFingerprint,
		keyFingerprint:          keyFingerprint,
		accountName:             accountName,

		hashFunc:   crypto.SHA1,
		privateKey: rsakey,
	}, nil
}

func (s *PrivateKeySigner) KeyFingerprint() string {
	return s.formattedKeyFingerprint
}

func (s *PrivateKeySigner) DefaultAlgorithm() string {
	return "rsa-sha1"
}

func (s *PrivateKeySigner) Sign(dateHeader string) (string, error) {
	const headerName = "date"

	hash := s.hashFunc.New()
	hash.Write([]byte(fmt.Sprintf("%s: %s", headerName, dateHeader)))
	digest := hash.Sum(nil)

	signed, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, s.hashFunc, digest)
	if err != nil {
		return "", errwrap.Wrapf("Error signing date header: {{err}}", err)
	}
	signedBase64 := base64.StdEncoding.EncodeToString(signed)

	return fmt.Sprintf(authorizationHeaderFormat, s.formattedKeyFingerprint, "rsa-sha1", headerName, signedBase64), nil
}

func (s *PrivateKeySigner) SignRaw(toSign string) (string, string, error) {
	hash := s.hashFunc.New()
	hash.Write([]byte(toSign))
	digest := hash.Sum(nil)

	signed, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, s.hashFunc, digest)
	if err != nil {
		return "", "", errwrap.Wrapf("Error signing string: {{err}}", err)
	}
	signedBase64 := base64.StdEncoding.EncodeToString(signed)

	return signedBase64, "rsa-sha1", nil
}
