package api

const APIVersion = "1"

// WellKnown is served at GET /.well-known/faramesh.json.
type WellKnown struct {
	APIVersion  string            `json:"api_version"`
	RegistryID  string            `json:"registry_id"`
	Search      string            `json:"search"`
	Artifact    ArtifactEndpoints `json:"artifact"`
	LegacyPacks string            `json:"legacy_packs_path,omitempty"`
	Trust       *TrustEndpoints   `json:"trust,omitempty"`
	Telemetry   *TelemetryEndpoints `json:"telemetry,omitempty"`
}

type TelemetryEndpoints struct {
	Stats string `json:"downloads_aggregate,omitempty"`
}

type ArtifactEndpoints struct {
	Providers  string `json:"providers"`
	Policies   string `json:"policies"`
	Frameworks string `json:"frameworks"`
	Versions   string `json:"versions_list,omitempty"`
}

type TrustEndpoints struct {
	OfficialKeyIDs []string `json:"official_key_ids,omitempty"`
	KeysURL        string   `json:"keys_url,omitempty"`
}

type SearchResponse struct {
	APIVersion string        `json:"api_version"`
	Packs      []PackSummary `json:"packs"`
}

type PackSummary struct {
	Kind          string `json:"kind"`
	Name          string `json:"name"`
	LatestVersion string `json:"latest_version"`
	Description   string `json:"description"`
	Downloads     int64  `json:"downloads,omitempty"`
	TrustTier     string `json:"trust_tier,omitempty"`
	Category      string `json:"category,omitempty"`
}

type VersionsResponse struct {
	APIVersion string         `json:"api_version"`
	Kind       string         `json:"kind"`
	Name       string         `json:"name"`
	Versions   []VersionEntry `json:"versions"`
}

type VersionEntry struct {
	Version   string `json:"version"`
	Changelog string `json:"changelog,omitempty"`
	Yanked    bool   `json:"yanked,omitempty"`
}

type StatsResponse struct {
	APIVersion      string `json:"api_version"`
	TotalArtifacts  int    `json:"total_artifacts"`
	OfficialCount   int    `json:"official_count"`
	ProviderCount   int    `json:"provider_count"`
	PolicyCount     int    `json:"policy_count"`
	FrameworkCount  int    `json:"framework_count"`
}

type PackVersionResponse struct {
	APIVersion  string `json:"api_version"`
	Kind        string `json:"kind,omitempty"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
	// PolicyFPL is canonical policy/framework bytes (FPL is the primary language).
	PolicyFPL string `json:"policy_fpl"`
	// PolicyYAML and PolicyJSON are optional sidecars when present in the catalog.
	PolicyYAML string `json:"policy_yaml,omitempty"`
	PolicyJSON string `json:"policy_json,omitempty"`
	SHA256Hex  string `json:"sha256_hex"`
	TrustTier  string `json:"trust_tier,omitempty"`
	Changelog  string `json:"changelog,omitempty"`
	ReadmeMarkdown string `json:"readme_markdown,omitempty"`

	ParameterSchema map[string]any   `json:"parameter_schema,omitempty"`
	RulesSummary    *RulesSummary     `json:"rules_summary,omitempty"`
	Compatibility   map[string]string `json:"compatibility,omitempty"`
	FarameshVersion string            `json:"faramesh_version,omitempty"`
	Signature       *ArtifactSignature `json:"signature,omitempty"`

	// Legacy hub field — omitted when policy_fpl is set; populated only for old YAML-only packs.
	PolicyYAMLLegacy string `json:"policy_yaml_legacy,omitempty"`
}

type RulesSummary struct {
	Permit []string `json:"permit,omitempty"`
	Deny   []string `json:"deny,omitempty"`
	Defer  []string `json:"defer,omitempty"`
}

type ArtifactSignature struct {
	Algorithm    string `json:"algorithm"`
	KeyID        string `json:"key_id,omitempty"`
	PublicKeyPEM string `json:"public_key_pem,omitempty"`
	PublicKeyB64 string `json:"public_key_b64,omitempty"`
	ValueB64     string `json:"value_b64"`
}

type ProviderDownload struct {
	URL       string `json:"url"`
	SHA256Hex string `json:"sha256_hex"`
	Size      int64  `json:"size,omitempty"`
}

type ProviderVersionResponse struct {
	APIVersion   string                      `json:"api_version"`
	Kind         string                      `json:"kind"`
	Name         string                      `json:"name"`
	Version      string                      `json:"version"`
	TrustTier    string                      `json:"trust_tier,omitempty"`
	DevOnly      bool                        `json:"dev_only,omitempty"`
	Capabilities []string                    `json:"capabilities,omitempty"`
	Downloads    map[string]ProviderDownload `json:"downloads"`
	ReadmeMarkdown string                    `json:"readme_markdown,omitempty"`
	ConfigSchema map[string]any              `json:"config_schema,omitempty"`
	Signature    *ArtifactSignature          `json:"signature,omitempty"`
}

type TrustKeysResponse struct {
	APIVersion string     `json:"api_version"`
	Keys       []TrustKey `json:"keys"`
}

type TrustKey struct {
	KeyID        string `json:"key_id"`
	Algorithm    string `json:"algorithm"`
	PublicKeyB64 string `json:"public_key_b64,omitempty"`
	PublicKeyPEM string `json:"public_key_pem,omitempty"`
	Purpose      string `json:"purpose,omitempty"`
}

type PublishRequest struct {
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
	// PolicyFPL is required for policy and framework publishes (FPL primary).
	PolicyFPL   string `json:"policy_fpl,omitempty"`
	PolicyYAML  string `json:"policy_yaml,omitempty"`
	PolicyJSON  string `json:"policy_json,omitempty"`
	ReadmeMarkdown string `json:"readme_markdown,omitempty"`
	TrustTier   string `json:"trust_tier,omitempty"`
	KeyID       string `json:"key_id,omitempty"`
	SignatureB64 string `json:"signature_b64,omitempty"`
}

type PublishResponse struct {
	APIVersion string `json:"api_version"`
	Accepted   bool   `json:"accepted"`
	Message    string `json:"message,omitempty"`
	ArtifactPath string `json:"artifact_path,omitempty"`
}
