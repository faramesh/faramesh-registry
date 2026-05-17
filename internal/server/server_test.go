package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestWellKnownAndSearch(t *testing.T) {
	root := filepath.Join("..", "..", "catalog")
	srv, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/.well-known/faramesh.json")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("well-known: %d", res.StatusCode)
	}

	res, err = http.Get(ts.URL + "/v1/search?kind=policy&q=stripe")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var body struct {
		Packs []struct {
			Kind string `json:"kind"`
			Name string `json:"name"`
		} `json:"packs"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range body.Packs {
		if p.Name == "faramesh/stripe" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected stripe in search, got %+v", body.Packs)
	}
}

func TestPolicyFPLPrimary(t *testing.T) {
	root := filepath.Join("..", "..", "catalog")
	srv, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/v1/policies/faramesh/demo/versions/0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var pv struct {
		PolicyFPL  string `json:"policy_fpl"`
		PolicyYAML string `json:"policy_yaml"`
	}
	if err := json.NewDecoder(res.Body).Decode(&pv); err != nil {
		t.Fatal(err)
	}
	if pv.PolicyFPL == "" {
		t.Fatal("policy_fpl required")
	}
	if pv.PolicyYAML == "" {
		t.Fatal("expected optional yaml sidecar")
	}
}
