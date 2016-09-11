package authentication

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/errwrap"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"strings"
)

type SSHAgentSigner struct {
	formattedKeyFingerprint string
	keyFingerprint          string
	accountName             string

	agent agent.Agent
	key   ssh.PublicKey
}

func NewSSHAgentSigner(keyFingerprint, accountName string) (*SSHAgentSigner, error) {
	sshAgentAddress := os.Getenv("SSH_AUTH_SOCK")
	if sshAgentAddress == "" {
		return nil, fmt.Errorf("SSH_AUTH_SOCK is not set")
	}

	conn, err := net.Dial("unix", sshAgentAddress)
	if err != nil {
		return nil, errwrap.Wrapf("Error dialing SSH agent: {{err}}", err)
	}

	ag := agent.NewClient(conn)

	keys, err := ag.List()
	if err != nil {
		return nil, errwrap.Wrapf("Error listing keys in SSH Agent: %s", err)
	}

	keyFingerprintMD5 := strings.Replace(keyFingerprint, ":", "", -1)

	var matchingKey ssh.PublicKey
	for _, key := range keys {
		h := md5.New()
		h.Write(key.Marshal())
		fp := fmt.Sprintf("%x", h.Sum(nil))

		if fp == keyFingerprintMD5 {
			matchingKey = key
		}
	}

	if matchingKey == nil {
		return nil, fmt.Errorf("No key in the SSH Agent matches fingerprint: %s", keyFingerprint)
	}

	return &SSHAgentSigner{
		formattedKeyFingerprint: formatPublicKeyFingerprint(matchingKey, true),
		keyFingerprint:          keyFingerprint,
		accountName:             accountName,
		agent:                   ag,
		key:                     matchingKey,
	}, nil
}

func (s *SSHAgentSigner) Sign(dateHeader string) (string, error) {
	const headerName = "date"

	signature, err := s.agent.Sign(s.key, []byte(fmt.Sprintf("%s: %s", headerName, dateHeader)))
	if err != nil {
		return "", errwrap.Wrapf("Error signing date header: {{err}}", err)
	}
	signedBase64 := base64.StdEncoding.EncodeToString(signature.Blob)

	var algorithm string
	switch signature.Format {
	case "ssh-rsa":
		algorithm = "rsa-sha1"
	default:
		return "", fmt.Errorf("Unsupported algorithm from SSH agent: %s", signature.Format)
	}

	return fmt.Sprintf(authorizationHeaderFormat, "/hashicorp/keys/jen20", algorithm, headerName, signedBase64), nil
}
