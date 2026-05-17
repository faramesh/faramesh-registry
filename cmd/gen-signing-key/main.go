// Command gen-signing-key generates an Ed25519 registry signing keypair for operators.
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	privB64 := base64.StdEncoding.EncodeToString(priv)
	pubB64 := base64.StdEncoding.EncodeToString(pub)

	fmt.Printf("REGISTRY_SIGNING_KEY_B64=%s\n", privB64)
	fmt.Println()
	fmt.Println("Public key entry for catalog/trust/keys.json:")
	entry := map[string]any{
		"id":             "faramesh-ed25519-2026",
		"algorithm":      "ed25519",
		"public_key_b64": pubB64,
		"owner":          "faramesh",
		"created_at":     time.Now().UTC().Format(time.RFC3339),
	}
	b, _ := json.MarshalIndent(map[string]any{"keys": []any{entry}}, "", "  ")
	fmt.Println(string(b))
}
