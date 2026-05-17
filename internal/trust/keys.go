package trust

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Store loads publisher public keys from catalog/trust/keys.json.
type Store struct {
	keys map[string]Key
}

type file struct {
	Keys []Key `json:"keys"`
}

// Key is a registry trust anchor.
type Key struct {
	ID           string `json:"id"`
	KeyID        string `json:"key_id"`
	Algorithm    string `json:"algorithm"`
	PublicKeyB64 string `json:"public_key_b64"`
	Owner        string `json:"owner,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Purpose      string `json:"purpose,omitempty"`
}

func Load(catalogDir string) (*Store, error) {
	path := filepath.Join(catalogDir, "trust", "keys.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return &Store{keys: map[string]Key{}}, nil
	}
	var f file
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, err
	}
	m := make(map[string]Key, len(f.Keys))
	for _, k := range f.Keys {
		if strings.TrimSpace(k.KeyID) == "" {
			k.KeyID = strings.TrimSpace(k.ID)
		}
		if k.KeyID == "" {
			continue
		}
		m[k.KeyID] = k
	}
	return &Store{keys: m}, nil
}

func (s *Store) All() []Key {
	if s == nil {
		return nil
	}
	out := make([]Key, 0, len(s.keys))
	for _, k := range s.keys {
		out = append(out, k)
	}
	return out
}

func (s *Store) Get(id string) (Key, error) {
	if s == nil {
		return Key{}, fmt.Errorf("trust store not loaded")
	}
	k, ok := s.keys[id]
	if !ok {
		return Key{}, fmt.Errorf("unknown key_id %q", id)
	}
	return k, nil
}
