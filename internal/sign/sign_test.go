package sign

import (
	"crypto/ed25519"
	"encoding/base64"
	"testing"
)

func TestSignVerifyRoundTrip(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	pub := priv.Public().(ed25519.PublicKey)
	payload := []byte(`agent "x" { default deny }`)
	sig, err := Sign(priv, payload)
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(payload, base64.StdEncoding.EncodeToString(pub), sig); err != nil {
		t.Fatal(err)
	}
}
