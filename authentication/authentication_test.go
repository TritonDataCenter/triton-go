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
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
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

// helper: generate an RSA private key and return the raw key, its PEM
// encoding, and the colon-separated MD5 fingerprint.
func generateRSAKey(t *testing.T) (*rsa.PrivateKey, []byte, string) {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	fingerprint := computeFingerprint(t, priv.Public())

	return priv, pemBytes, fingerprint
}

// helper: generate an ECDSA private key and return the raw key, its PEM
// encoding, and the colon-separated MD5 fingerprint.
func generateECDSAKey(t *testing.T) (*ecdsa.PrivateKey, []byte, string) {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey: %v", err)
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
// Test Group 1: ed25519 signature type
// ---------------------------------------------------------------------------

func TestEd25519Signature(t *testing.T) {
	blob := make([]byte, ed25519.SignatureSize)
	for i := range blob {
		blob[i] = byte(i)
	}

	sig, err := newEd25519Signature(blob)
	if err != nil {
		t.Fatalf("newEd25519Signature: %v", err)
	}

	if sig.SignatureType() != ED25519_SHA512 {
		t.Errorf("SignatureType() = %q, want %q", sig.SignatureType(), ED25519_SHA512)
	}

	expected := base64.StdEncoding.EncodeToString(blob)
	if sig.String() != expected {
		t.Errorf("String() = %q, want %q", sig.String(), expected)
	}
}

func TestEd25519SignatureValidation(t *testing.T) {
	tests := []struct {
		name string
		blob []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
		{"too short", make([]byte, 32)},
		{"too long", make([]byte, 128)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newEd25519Signature(tc.blob)
			if err == nil {
				t.Error("expected error for invalid signature length, got nil")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Test Group 2: keyFormatToKeyType
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
// Test Group 3: formatPublicKeyFingerprint
// ---------------------------------------------------------------------------

func TestFormatPublicKeyFingerprint(t *testing.T) {
	rsaKey, _, rsaFP := generateRSAKey(t)
	edKey, _, edFP := generateEd25519Key(t)

	// Test ed25519 private key (the new code path)
	t.Run("ed25519 private key", func(t *testing.T) {
		got, err := formatPublicKeyFingerprint(edKey, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != edFP {
			t.Errorf("fingerprint = %q, want %q", got, edFP)
		}
	})

	// Test ed25519 raw hex (display=false)
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

	// Test RSA private key
	t.Run("rsa private key", func(t *testing.T) {
		got, err := formatPublicKeyFingerprint(rsaKey, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != rsaFP {
			t.Errorf("fingerprint = %q, want %q", got, rsaFP)
		}
	})

	// Test ecdsa private key
	t.Run("ecdsa private key", func(t *testing.T) {
		ecKey, _, ecFP := generateECDSAKey(t)
		got, err := formatPublicKeyFingerprint(ecKey, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != ecFP {
			t.Errorf("fingerprint = %q, want %q", got, ecFP)
		}
	})

	// Test ssh.PublicKey directly (the new case added in this branch)
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

	// Test unsupported type
	t.Run("unsupported type", func(t *testing.T) {
		_, err := formatPublicKeyFingerprint("not a key", false)
		if err == nil {
			t.Error("expected error for unsupported type, got nil")
		}
	})
}

// ---------------------------------------------------------------------------
// Test Group 4: PrivateKeySigner
// ---------------------------------------------------------------------------

func TestPrivateKeySignerEd25519(t *testing.T) {
	_, pemBytes, fingerprint := generateEd25519Key(t)

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

		if !strings.Contains(header, "ed25519-sha512") {
			t.Errorf("Sign output missing algorithm, got: %s", header)
		}
		if !strings.Contains(header, "Signature keyId=") {
			t.Errorf("Sign output missing keyId, got: %s", header)
		}
		if !strings.Contains(header, "testaccount") {
			t.Errorf("Sign output missing account name, got: %s", header)
		}
	})

	t.Run("SignRaw", func(t *testing.T) {
		signedBase64, algo, err := signer.SignRaw("test message")
		if err != nil {
			t.Fatalf("SignRaw: %v", err)
		}

		if algo != ED25519_SHA512 {
			t.Errorf("algorithm = %q, want %q", algo, ED25519_SHA512)
		}

		// Verify the signature is valid base64
		_, err = base64.StdEncoding.DecodeString(signedBase64)
		if err != nil {
			t.Fatalf("invalid base64 in signature: %v", err)
		}
	})

	t.Run("SignRaw verify", func(t *testing.T) {
		// Generate a fresh key so we have access to the raw private key
		// for verification.
		priv, pemBytes2, fp2 := generateEd25519Key(t)

		signer2, err := NewPrivateKeySigner(PrivateKeySignerInput{
			KeyID:              fp2,
			PrivateKeyMaterial: pemBytes2,
			AccountName:        "testaccount",
		})
		if err != nil {
			t.Fatalf("NewPrivateKeySigner: %v", err)
		}

		message := "verify me"
		signedBase64, _, err := signer2.SignRaw(message)
		if err != nil {
			t.Fatalf("SignRaw: %v", err)
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

func TestPrivateKeySignerRSA(t *testing.T) {
	_, pemBytes, fingerprint := generateRSAKey(t)

	signer, err := NewPrivateKeySigner(PrivateKeySignerInput{
		KeyID:              fingerprint,
		PrivateKeyMaterial: pemBytes,
		AccountName:        "testaccount",
	})
	if err != nil {
		t.Fatalf("NewPrivateKeySigner: %v", err)
	}

	if signer.DefaultAlgorithm() != RSA_SHA512 {
		t.Errorf("DefaultAlgorithm() = %q, want %q", signer.DefaultAlgorithm(), RSA_SHA512)
	}

	header, err := signer.Sign("Thu, 01 Jan 2025 00:00:00 GMT", false)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if !strings.Contains(header, "rsa-sha512") {
		t.Errorf("Sign output missing algorithm, got: %s", header)
	}
}

func TestPrivateKeySignerECDSA(t *testing.T) {
	_, pemBytes, fingerprint := generateECDSAKey(t)

	signer, err := NewPrivateKeySigner(PrivateKeySignerInput{
		KeyID:              fingerprint,
		PrivateKeyMaterial: pemBytes,
		AccountName:        "testaccount",
	})
	if err != nil {
		t.Fatalf("NewPrivateKeySigner: %v", err)
	}

	if signer.DefaultAlgorithm() != ECDSA_SHA512 {
		t.Errorf("DefaultAlgorithm() = %q, want %q", signer.DefaultAlgorithm(), ECDSA_SHA512)
	}

	header, err := signer.Sign("Thu, 01 Jan 2025 00:00:00 GMT", false)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if !strings.Contains(header, "ecdsa-sha512") {
		t.Errorf("Sign output missing algorithm, got: %s", header)
	}
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

func TestPrivateKeySignerBadKeyMaterial(t *testing.T) {
	_, err := NewPrivateKeySigner(PrivateKeySignerInput{
		KeyID:              "00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00",
		PrivateKeyMaterial: []byte("not a valid key"),
		AccountName:        "testaccount",
	})
	if err == nil {
		t.Fatal("expected error for bad key material, got nil")
	}
}

// Test that Sign() produces a consistent format by calling SignRaw() internally.
func TestPrivateKeySignerSignConsistency(t *testing.T) {
	_, pemBytes, fingerprint := generateEd25519Key(t)

	signer, err := NewPrivateKeySigner(PrivateKeySignerInput{
		KeyID:              fingerprint,
		PrivateKeyMaterial: pemBytes,
		AccountName:        "testaccount",
		Username:           "testuser",
	})
	if err != nil {
		t.Fatalf("NewPrivateKeySigner: %v", err)
	}

	// Sign should produce a properly formatted authorization header
	header, err := signer.Sign("Thu, 01 Jan 2025 00:00:00 GMT", false)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	// Verify the header format
	if !strings.HasPrefix(header, "Signature keyId=\"") {
		t.Errorf("header should start with 'Signature keyId=\"', got: %s", header)
	}
	if !strings.Contains(header, `algorithm="ed25519-sha512"`) {
		t.Errorf("header missing algorithm, got: %s", header)
	}
	if !strings.Contains(header, `headers="date"`) {
		t.Errorf("header missing headers field, got: %s", header)
	}
	if !strings.Contains(header, `signature="`) {
		t.Errorf("header missing signature field, got: %s", header)
	}
}

// ---------------------------------------------------------------------------
// Test Group 5: KeyID generation
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
