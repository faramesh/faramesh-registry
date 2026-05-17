// Command sign-catalog signs canonical .fpl artifacts in the catalog (GitOps CI).
package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/faramesh/faramesh-registry/internal/catalog"
	"github.com/faramesh/faramesh-registry/internal/sign"
)

func main() {
	catalogDir := flag.String("catalog", "catalog", "catalog root")
	keyB64 := flag.String("key", os.Getenv("REGISTRY_SIGNING_KEY_B64"), "Ed25519 private key (64-byte seed+pub, base64)")
	keyID := flag.String("key-id", "faramesh-ed25519-2026", "signing key id")
	providers := flag.Bool("providers", false, "also sign provider binaries and update manifest.json signatures")
	flag.Parse()
	keyEnv := strings.TrimSpace(os.Getenv("REGISTRY_SIGNING_KEY_B64"))
	if strings.TrimSpace(*keyB64) == "" {
		*keyB64 = keyEnv
	}
	if strings.TrimSpace(*keyB64) == "" {
		fmt.Fprintln(os.Stderr, "REGISTRY_SIGNING_KEY_B64 is not set (export the operator private key from gen-signing-key)")
		os.Exit(1)
	}
	privRaw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(*keyB64))
	if err != nil {
		fatal(err)
	}
	if len(privRaw) != ed25519.PrivateKeySize {
		fatal(fmt.Errorf("private key must be %d bytes, got %d", ed25519.PrivateKeySize, len(privRaw)))
	}
	priv := ed25519.PrivateKey(privRaw)
	pub := priv.Public().(ed25519.PublicKey)

	idx, base, err := catalog.LoadIndex(*catalogDir)
	if err != nil {
		fatal(err)
	}
	n := 0
	signFPL := func(path, primary string) {
		artDir := filepath.Dir(path)
		if st, err := os.Stat(path); err == nil && st.IsDir() {
			artDir = path
		}
		fplPath := filepath.Join(artDir, primary)
		body, err := os.ReadFile(fplPath)
		if err != nil {
			return
		}
		sig, err := sign.Sign(priv, body)
		if err != nil {
			fatal(err)
		}
		ref := catalog.SignatureRef{
			KeyID: *keyID, Algorithm: "ed25519",
			ValueB64: sig, PublicKeyB64: base64.StdEncoding.EncodeToString(pub),
		}
		b, _ := json.MarshalIndent(ref, "", "  ")
		if err := os.WriteFile(fplPath+".sig", b, 0o644); err != nil {
			fatal(err)
		}
		n++
		fmt.Println("signed", fplPath)
	}
	for _, e := range idx.Policies {
		for _, p := range e.Versions {
			signFPL(p, "policy.fpl")
		}
	}
	for _, e := range idx.Frameworks {
		for _, p := range e.Versions {
			signFPL(p, "profile.fpl")
		}
	}
	for _, e := range idx.Packs {
		for _, p := range e.Versions {
			signFPL(p, "policy.fpl")
		}
	}
	if *providers {
		for _, e := range idx.Providers {
			for _, manifestPath := range e.Versions {
				if err := signProviderManifest(*catalogDir, manifestPath, priv, *keyID, pub); err != nil {
					fatal(err)
				}
				n++
			}
		}
	}
	_ = base
	fmt.Printf("signed %d artifacts\n", n)
}

func signProviderManifest(catalogDir, manifestRel string, priv ed25519.PrivateKey, keyID string, pub ed25519.PublicKey) error {
	manifestPath := manifestRel
	if !filepath.IsAbs(manifestPath) {
		manifestPath = filepath.Join(catalogDir, manifestRel)
	}
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	var man catalog.ProviderManifest
	if err := json.Unmarshal(b, &man); err != nil {
		return err
	}
	artDir := filepath.Dir(manifestPath)
	for plat, dl := range man.Downloads {
		binPath := strings.TrimPrefix(strings.TrimSpace(dl.URL), "file://")
		if binPath == "" {
			continue
		}
		if !filepath.IsAbs(binPath) {
			binPath = filepath.Join(artDir, binPath)
		}
		body, err := os.ReadFile(binPath)
		if err != nil {
			return fmt.Errorf("%s: %w", plat, err)
		}
		sig, err := sign.Sign(priv, body)
		if err != nil {
			return err
		}
		if err := os.WriteFile(binPath+".sig", []byte(sig), 0o644); err != nil {
			return err
		}
		sum := sha256.Sum256(body)
		man.Downloads[plat] = catalog.ProviderDownload{
			URL:       dl.URL,
			SHA256Hex: hex.EncodeToString(sum[:]),
			Size:      int64(len(body)),
		}
		fmt.Println("signed", binPath)
	}
	man.Signature = nil
	canonical, err := json.Marshal(man)
	if err != nil {
		return err
	}
	manifestSig, err := sign.Sign(priv, canonical)
	if err != nil {
		return err
	}
	man.Signature = &catalog.SignatureRef{
		KeyID: keyID, Algorithm: "ed25519",
		ValueB64: manifestSig, PublicKeyB64: base64.StdEncoding.EncodeToString(pub),
	}
	out, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(manifestPath, append(out, '\n'), 0o644)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
