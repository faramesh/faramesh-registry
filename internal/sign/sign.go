// Package sign verifies Ed25519 signatures over registry artifact bytes.
package sign

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// Verify checks a detached Ed25519 signature over payload.
func Verify(payload []byte, publicKeyB64, signatureB64 string) error {
	pub, err := decodePublicKey(publicKeyB64)
	if err != nil {
		return err
	}
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(signatureB64))
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	if len(sig) != ed25519.SignatureSize {
		return fmt.Errorf("invalid signature length %d", len(sig))
	}
	if !ed25519.Verify(pub, payload, sig) {
		return errors.New("signature mismatch")
	}
	return nil
}

func decodePublicKey(b64 string) (ed25519.PublicKey, error) {
	b64 = strings.TrimSpace(b64)
	if b64 == "" {
		return nil, errors.New("empty public key")
	}
	if strings.HasPrefix(b64, "-----BEGIN") {
		block, _ := pem.Decode([]byte(b64))
		if block == nil {
			return nil, errors.New("invalid PEM public key")
		}
		if len(block.Bytes) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("unexpected PEM key size %d", len(block.Bytes))
		}
		return ed25519.PublicKey(block.Bytes), nil
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("unexpected public key length %d", len(raw))
	}
	return ed25519.PublicKey(raw), nil
}

// Sign creates a detached Ed25519 signature.
func Sign(privateKey ed25519.PrivateKey, payload []byte) (string, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return "", errors.New("invalid private key")
	}
	sig := ed25519.Sign(privateKey, payload)
	return base64.StdEncoding.EncodeToString(sig), nil
}
