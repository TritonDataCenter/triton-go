//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

//
// Copyright 2025 Edgecast Cloud LLC.
//

package authentication

import (
	"crypto"
	"crypto/ed25519"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

// helper: generate an ed25519 private key and return the raw key, its PEM
// encoding, and the colon-separated MD5 fingerprint suitable for use as a
// KeyID.
func generateEd25519Key(t *testing.T) (ed25519.PrivateKey, []byte, string) {
	t.Helper()

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("ed25519.GenerateKey: %v", err)
	}

	pkcs8, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatalf("MarshalPKCS8PrivateKey: %v", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8})

	fingerprint := computeFingerprint(t, priv.Public())

	return priv, pemBytes, fingerprint
}

// computeFingerprint returns the colon-separated MD5 fingerprint of a
// crypto.PublicKey, matching the format expected by NewPrivateKeySigner.
func computeFingerprint(t *testing.T, pub crypto.PublicKey) string {
	t.Helper()

	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatalf("ssh.NewPublicKey: %v", err)
	}

	hash := md5.New()
	hash.Write(sshPub.Marshal())
	hex := fmt.Sprintf("%x", hash.Sum(nil))

	var parts []string
	for i := 0; i < len(hex); i += 2 {
		parts = append(parts, hex[i:i+2])
	}

	return strings.Join(parts, ":")
}

// ---------------------------------------------------------------------------
// ed25519 signature validation
// ---------------------------------------------------------------------------

func TestEd25519SignatureValidation(t *testing.T) {
	_, err := newEd25519Signature(make([]byte, 32))
	if err == nil {
		t.Error("expected error for invalid signature length, got nil")
	}
}

// ---------------------------------------------------------------------------
// keyFormatToKeyType
// ---------------------------------------------------------------------------

func TestKeyFormatToKeyType(t *testing.T) {
	tests := []struct {
		format  string
		want    string
		wantErr bool
	}{
		{"ssh-rsa", "rsa", false},
		{"ssh-ed25519", "ed25519", false},
		{"ecdsa-sha2-nistp256", "ecdsa", false},
		{"ecdsa-sha2-nistp384", "ecdsa", false},
		{"ecdsa-sha2-nistp521", "ecdsa", false},
		{"ssh-dss", "", true},
		{"unknown", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.format, func(t *testing.T) {
			got, err := keyFormatToKeyType(tc.format)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for format %q, got nil", tc.format)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("keyFormatToKeyType(%q) = %q, want %q", tc.format, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// formatPublicKeyFingerprint — only new code paths
// ---------------------------------------------------------------------------

func TestFormatPublicKeyFingerprint(t *testing.T) {
	edKey, _, edFP := generateEd25519Key(t)

	t.Run("ed25519 private key", func(t *testing.T) {
		got, err := formatPublicKeyFingerprint(edKey, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != edFP {
			t.Errorf("fingerprint = %q, want %q", got, edFP)
		}
	})

	t.Run("ed25519 raw hex", func(t *testing.T) {
		got, err := formatPublicKeyFingerprint(edKey, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := strings.Replace(edFP, ":", "", -1)
		if got != expected {
			t.Errorf("raw fingerprint = %q, want %q", got, expected)
		}
	})

	t.Run("ssh.PublicKey", func(t *testing.T) {
		sshPub, err := ssh.NewPublicKey(edKey.Public())
		if err != nil {
			t.Fatalf("ssh.NewPublicKey: %v", err)
		}
		got, err := formatPublicKeyFingerprint(sshPub, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != edFP {
			t.Errorf("fingerprint = %q, want %q", got, edFP)
		}
	})

	t.Run("unsupported type", func(t *testing.T) {
		_, err := formatPublicKeyFingerprint("not a key", false)
		if err == nil {
			t.Error("expected error for unsupported type, got nil")
		}
	})
}

// ---------------------------------------------------------------------------
// PrivateKeySigner with ed25519
// ---------------------------------------------------------------------------

func TestPrivateKeySignerEd25519(t *testing.T) {
	priv, pemBytes, fingerprint := generateEd25519Key(t)

	signer, err := NewPrivateKeySigner(PrivateKeySignerInput{
		KeyID:              fingerprint,
		PrivateKeyMaterial: pemBytes,
		AccountName:        "testaccount",
		Username:           "testuser",
	})
	if err != nil {
		t.Fatalf("NewPrivateKeySigner: %v", err)
	}

	if signer.DefaultAlgorithm() != ED25519_SHA512 {
		t.Errorf("DefaultAlgorithm() = %q, want %q", signer.DefaultAlgorithm(), ED25519_SHA512)
	}

	t.Run("Sign", func(t *testing.T) {
		header, err := signer.Sign("Thu, 01 Jan 2025 00:00:00 GMT", false)
		if err != nil {
			t.Fatalf("Sign: %v", err)
		}

		if !strings.HasPrefix(header, "Signature keyId=\"") {
			t.Errorf("header should start with 'Signature keyId=\"', got: %s", header)
		}
		if !strings.Contains(header, `algorithm="ed25519-sha512"`) {
			t.Errorf("header missing algorithm, got: %s", header)
		}
		if !strings.Contains(header, `headers="date"`) {
			t.Errorf("header missing headers field, got: %s", header)
		}
		if !strings.Contains(header, "testaccount") {
			t.Errorf("header missing account name, got: %s", header)
		}
	})

	t.Run("SignRaw verify", func(t *testing.T) {
		message := "verify me"
		signedBase64, algo, err := signer.SignRaw(message)
		if err != nil {
			t.Fatalf("SignRaw: %v", err)
		}

		if algo != ED25519_SHA512 {
			t.Errorf("algorithm = %q, want %q", algo, ED25519_SHA512)
		}

		sigBytes, err := base64.StdEncoding.DecodeString(signedBase64)
		if err != nil {
			t.Fatalf("base64 decode: %v", err)
		}

		pub := priv.Public().(ed25519.PublicKey)
		if !ed25519.Verify(pub, []byte(message), sigBytes) {
			t.Error("ed25519.Verify failed: signature is not valid")
		}
	})
}

func TestPrivateKeySignerBadFingerprint(t *testing.T) {
	_, pemBytes, _ := generateEd25519Key(t)

	_, err := NewPrivateKeySigner(PrivateKeySignerInput{
		KeyID:              "00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00",
		PrivateKeyMaterial: pemBytes,
		AccountName:        "testaccount",
	})
	if err == nil {
		t.Fatal("expected error for mismatched fingerprint, got nil")
	}
	if !strings.Contains(err.Error(), "does not match") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// ---------------------------------------------------------------------------
// KeyID generation
// ---------------------------------------------------------------------------

func TestKeyIDGenerate(t *testing.T) {
	tests := []struct {
		name   string
		input  KeyID
		expect string
	}{
		{
			name: "no user",
			input: KeyID{
				AccountName: "myaccount",
				Fingerprint: "aa:bb:cc",
			},
			expect: "/myaccount/keys/aa:bb:cc",
		},
		{
			name: "with user",
			input: KeyID{
				AccountName: "myaccount",
				UserName:    "myuser",
				Fingerprint: "aa:bb:cc",
			},
			expect: "/myaccount/users/myuser/keys/aa:bb:cc",
		},
		{
			name: "manta with user",
			input: KeyID{
				AccountName: "myaccount",
				UserName:    "myuser",
				Fingerprint: "aa:bb:cc",
				IsManta:     true,
			},
			expect: "/myaccount/myuser/keys/aa:bb:cc",
		},
		{
			name: "manta no user",
			input: KeyID{
				AccountName: "myaccount",
				Fingerprint: "aa:bb:cc",
				IsManta:     true,
			},
			expect: "/myaccount/keys/aa:bb:cc",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.generate()
			if got != tc.expect {
				t.Errorf("generate() = %q, want %q", got, tc.expect)
			}
		})
	}
}
