// Package catalog loads the GitOps registry index and artifact sidecars.
package catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Index is catalog.json at the catalog root.
type Index struct {
	RegistryID string           `json:"registry_id"`
	Packs      []Entry          `json:"packs"`
	Providers  []ProviderEntry  `json:"providers,omitempty"`
	Policies   []Entry          `json:"policies,omitempty"`
	Frameworks []Entry          `json:"frameworks,omitempty"`
}

type Entry struct {
	Name          string            `json:"name"`
	LatestVersion string            `json:"latest_version"`
	Description   string            `json:"description"`
	TrustTier     string            `json:"trust_tier,omitempty"`
	Category      string            `json:"category,omitempty"`
	Capabilities  []string          `json:"capabilities,omitempty"`
	Versions      map[string]string `json:"versions"`
}

type ProviderEntry struct {
	Entry
	Capabilities []string `json:"capabilities,omitempty"`
}

// VersionMeta is optional meta.json beside canonical FPL/manifest files.
type VersionMeta struct {
	Description     string            `json:"description,omitempty"`
	Changelog       string            `json:"changelog,omitempty"`
	ParameterSchema map[string]any      `json:"parameter_schema,omitempty"`
	RulesSummary    *RulesSummary       `json:"rules_summary,omitempty"`
	Compatibility   map[string]string   `json:"compatibility,omitempty"`
	FarameshVersion string            `json:"faramesh_version,omitempty"`
	ConfigSchema    map[string]any      `json:"config_schema,omitempty"`
	Signature       *SignatureRef       `json:"signature,omitempty"`
}

type RulesSummary struct {
	Permit []string `json:"permit,omitempty"`
	Deny   []string `json:"deny,omitempty"`
	Defer  []string `json:"defer,omitempty"`
}

type SignatureRef struct {
	KeyID        string `json:"key_id"`
	Algorithm    string `json:"algorithm"`
	ValueB64     string `json:"value_b64"`
	PublicKeyB64 string `json:"public_key_b64,omitempty"`
}

type ProviderManifest struct {
	Name         string                      `json:"name,omitempty"`
	Version      string                      `json:"version,omitempty"`
	DevOnly      bool                        `json:"dev_only,omitempty"`
	Capabilities []string                    `json:"capabilities,omitempty"`
	ConfigSchema map[string]any              `json:"config_schema,omitempty"`
	Downloads    map[string]ProviderDownload `json:"downloads"`
	Signature    *SignatureRef               `json:"signature,omitempty"`
}

type ProviderDownload struct {
	URL       string `json:"url"`
	SHA256Hex string `json:"sha256_hex"`
	Size      int64  `json:"size,omitempty"`
}

// FPLArtifact is the canonical FPL bytes plus optional YAML/JSON sidecars.
type FPLArtifact struct {
	FPLBytes    []byte
	YAMLBytes   []byte
	JSONBytes   []byte
	Readme      []byte
	Meta        VersionMeta
	SHA256Hex   string
	Signature   *SignatureRef
}

// LoadIndex reads catalog.json and resolves version paths to absolute paths.
func LoadIndex(catalogDir string) (*Index, string, error) {
	catalogDir, err := filepath.Abs(catalogDir)
	if err != nil {
		return nil, "", err
	}
	path := filepath.Join(catalogDir, "catalog.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("read catalog.json: %w", err)
	}
	var idx Index
	if err := json.Unmarshal(b, &idx); err != nil {
		return nil, "", err
	}
	resolve := func(m map[string]string) {
		for ver, rel := range m {
			if rel != "" && !filepath.IsAbs(rel) {
				m[ver] = filepath.Join(catalogDir, rel)
			}
		}
	}
	for i := range idx.Packs {
		resolve(idx.Packs[i].Versions)
	}
	for i := range idx.Providers {
		resolve(idx.Providers[i].Versions)
	}
	for i := range idx.Policies {
		resolve(idx.Policies[i].Versions)
	}
	for i := range idx.Frameworks {
		resolve(idx.Frameworks[i].Versions)
	}
	return &idx, catalogDir, nil
}

// LoadFPLArtifact loads policy.fpl or profile.fpl and optional sidecars from a version directory.
func LoadFPLArtifact(versionPath string, primaryName string) (*FPLArtifact, error) {
	dir := filepath.Dir(versionPath)
	if st, err := os.Stat(versionPath); err == nil && !st.IsDir() {
		dir = filepath.Dir(versionPath)
	} else if st != nil && st.IsDir() {
		dir = versionPath
		versionPath = filepath.Join(dir, primaryName)
	}
	fplPath := filepath.Join(dir, primaryName)
	fpl, err := os.ReadFile(fplPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", primaryName, err)
	}
	out := &FPLArtifact{FPLBytes: fpl}
	sum := sha256.Sum256(fpl)
	out.SHA256Hex = hex.EncodeToString(sum[:])

	if b, err := os.ReadFile(filepath.Join(dir, strings.Replace(primaryName, ".fpl", ".yaml", 1))); err == nil {
		out.YAMLBytes = b
	} else if b, err := os.ReadFile(filepath.Join(dir, "policy.yaml")); err == nil {
		out.YAMLBytes = b
	} else if b, err := os.ReadFile(filepath.Join(dir, "profile.yaml")); err == nil {
		out.YAMLBytes = b
	}
	if b, err := os.ReadFile(filepath.Join(dir, strings.Replace(primaryName, ".fpl", ".json", 1))); err == nil {
		out.JSONBytes = b
	} else if b, err := os.ReadFile(filepath.Join(dir, "policy.json")); err == nil {
		out.JSONBytes = b
	} else if b, err := os.ReadFile(filepath.Join(dir, "profile.json")); err == nil {
		out.JSONBytes = b
	}
	if b, err := os.ReadFile(filepath.Join(dir, "README.md")); err == nil {
		out.Readme = b
	}
	if b, err := os.ReadFile(filepath.Join(dir, "meta.json")); err == nil {
		_ = json.Unmarshal(b, &out.Meta)
	}
	if out.Meta.Signature == nil {
		if sig, err := readSigFile(dir, primaryName); err == nil {
			out.Signature = sig
		}
	} else {
		out.Signature = out.Meta.Signature
	}
	return out, nil
}

func readSigFile(dir, primaryName string) (*SignatureRef, error) {
	sigPath := filepath.Join(dir, primaryName+".sig")
	b, err := os.ReadFile(sigPath)
	if err != nil {
		return nil, err
	}
	var ref SignatureRef
	if err := json.Unmarshal(b, &ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

// LoadProviderManifest reads provider manifest.json.
func LoadProviderManifest(path string) (*ProviderManifest, []byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var man ProviderManifest
	if err := json.Unmarshal(b, &man); err != nil {
		return nil, b, err
	}
	return &man, b, nil
}
