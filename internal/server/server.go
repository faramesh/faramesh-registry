// Package server implements the Faramesh Registry v1 HTTP API.
// FPL is the canonical artifact language; YAML and JSON are optional sidecars.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/faramesh/faramesh-registry/internal/api"
	"github.com/faramesh/faramesh-registry/internal/catalog"
	"github.com/faramesh/faramesh-registry/internal/sign"
	"github.com/faramesh/faramesh-registry/internal/trust"
)

// Server serves registry HTTP routes.
type Server struct {
	baseDir    string
	index      *catalog.Index
	trust      *trust.Store
	allowWrite bool
}

// New loads catalog.json and trust keys from catalogDir.
func New(catalogDir string) (*Server, error) {
	idx, base, err := catalog.LoadIndex(catalogDir)
	if err != nil {
		return nil, err
	}
	ts, err := trust.Load(base)
	if err != nil {
		return nil, err
	}
	return &Server{
		baseDir:    base,
		index:      idx,
		trust:      ts,
		allowWrite: os.Getenv("REGISTRY_PUBLISH_WRITE") == "1",
	}, nil
}

// Handler returns the root HTTP handler.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/faramesh.json", s.handleWellKnown)
	mux.HandleFunc("GET /v1/search", s.handleSearch)
	mux.HandleFunc("GET /v1/stats", s.handleStats)
	mux.HandleFunc("GET /v1/trust/keys", s.handleTrustKeys)
	mux.HandleFunc("GET /v1/trust/keys/{key_id}", s.handleTrustKey)
	mux.HandleFunc("POST /v1/publish", s.handlePublish)
	mux.HandleFunc("GET /v1/packs/{name}/versions/{version}", s.handleLegacyPack)
	mux.HandleFunc("GET /v1/providers/", s.routeProviders)
	mux.HandleFunc("GET /v1/policies/", s.routePolicies)
	mux.HandleFunc("GET /v1/frameworks/", s.routeFrameworks)
	artDir := filepath.Join(s.baseDir, "artifacts")
	mux.Handle("/artifacts/", http.StripPrefix("/artifacts/", http.FileServer(http.Dir(artDir))))
	return corsMiddleware(mux)
}

func (s *Server) handleWellKnown(w http.ResponseWriter, _ *http.Request) {
	id := s.index.RegistryID
	if id == "" {
		id = "registry.faramesh.dev"
	}
	ids := []string{"faramesh-ed25519-dev"}
	for _, k := range s.trust.All() {
		ids = append(ids, k.KeyID)
	}
	writeJSON(w, api.WellKnown{
		APIVersion: APIVersion,
		RegistryID: id,
		Search:     "/v1/search",
		Artifact: api.ArtifactEndpoints{
			Providers:  "/v1/providers/{name}/versions/{version}",
			Policies:   "/v1/policies/{name}/versions/{version}",
			Frameworks: "/v1/frameworks/{name}/versions/{version}",
			Versions:   "/v1/policies/{name}/versions",
		},
		LegacyPacks: "/v1/packs/{name}/versions/{version}",
		Trust: &api.TrustEndpoints{
			OfficialKeyIDs: ids,
			KeysURL:        "/v1/trust/keys",
		},
		Telemetry: &api.TelemetryEndpoints{Stats: "/v1/stats"},
	})
}

func (s *Server) handleStats(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, api.StatsResponse{
		APIVersion:     APIVersion,
		ProviderCount:  len(s.index.Providers),
		PolicyCount:    len(s.index.Policies) + len(s.index.Packs),
		FrameworkCount: len(s.index.Frameworks),
		TotalArtifacts: len(s.index.Providers) + len(s.index.Policies) + len(s.index.Frameworks) + len(s.index.Packs),
		OfficialCount:  s.countTier("official"),
	})
}

func (s *Server) countTier(tier string) int {
	n := 0
	for _, e := range s.index.Providers {
		if e.TrustTier == tier {
			n++
		}
	}
	for _, e := range s.index.Policies {
		if e.TrustTier == tier {
			n++
		}
	}
	for _, e := range s.index.Frameworks {
		if e.TrustTier == tier {
			n++
		}
	}
	return n
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	kind := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("kind")))
	tier := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tier")))
	var rows []api.PackSummary
	add := func(name, ver, desc, t, k, cat string) {
		if kind != "" && kind != k {
			return
		}
		if tier != "" && strings.ToLower(t) != tier {
			return
		}
		if q != "" && !strings.Contains(strings.ToLower(name), q) && !strings.Contains(strings.ToLower(desc), q) {
			return
		}
		rows = append(rows, api.PackSummary{Kind: k, Name: name, LatestVersion: ver, Description: desc, TrustTier: t, Category: cat})
	}
	for _, p := range s.index.Providers {
		add(p.Name, p.LatestVersion, p.Description, p.TrustTier, "provider", strings.Join(p.Capabilities, ","))
	}
	for _, p := range s.index.Policies {
		add(p.Name, p.LatestVersion, p.Description, p.TrustTier, "policy", p.Category)
	}
	for _, p := range s.index.Frameworks {
		add(p.Name, p.LatestVersion, p.Description, p.TrustTier, "framework", p.Category)
	}
	for _, p := range s.index.Packs {
		add(p.Name, p.LatestVersion, p.Description, p.TrustTier, "policy", p.Category)
	}
	writeJSON(w, api.SearchResponse{APIVersion: APIVersion, Packs: rows})
}

func (s *Server) handleTrustKeys(w http.ResponseWriter, _ *http.Request) {
	var keys []api.TrustKey
	for _, k := range s.trust.All() {
		keys = append(keys, trustKeyAPI(k))
	}
	if len(keys) == 0 {
		keys = []api.TrustKey{{
			KeyID: "faramesh-ed25519-dev", Algorithm: "ed25519",
			PublicKeyB64: "IGchbOFuk05IUVcEEsOV4iSJ9S7CFv+QMOTY1dR2B4I=",
			Purpose: "artifact-signature",
		}}
	}
	writeJSON(w, api.TrustKeysResponse{APIVersion: APIVersion, Keys: keys})
}

func (s *Server) handleTrustKey(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("key_id"))
	k, err := s.trust.Get(id)
	if err != nil {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	writeJSON(w, trustKeyAPI(k))
}

func trustKeyAPI(k trust.Key) api.TrustKey {
	return api.TrustKey{
		KeyID: k.KeyID, Algorithm: k.Algorithm,
		PublicKeyB64: k.PublicKeyB64, Purpose: k.Purpose,
	}
}

func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request) {
	var req api.PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid publish payload", http.StatusBadRequest)
		return
	}
	req.Kind = strings.ToLower(strings.TrimSpace(req.Kind))
	req.Name = strings.TrimSpace(req.Name)
	req.Version = strings.TrimSpace(req.Version)
	if req.Kind == "" || req.Name == "" || req.Version == "" {
		http.Error(w, "kind, name, and version are required", http.StatusBadRequest)
		return
	}
	if req.Kind == "policy" || req.Kind == "framework" {
		if strings.TrimSpace(req.PolicyFPL) == "" {
			http.Error(w, "policy_fpl is required (FPL is the primary language)", http.StatusBadRequest)
			return
		}
	}
	if req.SignatureB64 != "" {
		keyID := req.KeyID
		if keyID == "" {
			keyID = "faramesh-ed25519-dev"
		}
		tk, err := s.trust.Get(keyID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		payload := []byte(req.PolicyFPL)
		if len(payload) == 0 {
			http.Error(w, "nothing to verify", http.StatusBadRequest)
			return
		}
		if err := sign.Verify(payload, tk.PublicKeyB64, req.SignatureB64); err != nil {
			http.Error(w, "signature verification failed: "+err.Error(), http.StatusUnauthorized)
			return
		}
	}
	if !s.allowWrite {
		writeJSON(w, api.PublishResponse{
			APIVersion: APIVersion, Accepted: true,
			Message: "signature OK; merge catalog via GitOps PR (set REGISTRY_PUBLISH_WRITE=1 for local dev writes)",
		})
		return
	}
	path, err := s.writePublish(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, api.PublishResponse{
		APIVersion: APIVersion, Accepted: true,
		Message:    "artifact written",
		ArtifactPath: path,
	})
}

func (s *Server) writePublish(req api.PublishRequest) (string, error) {
	primary := "policy.fpl"
	subdir := "policies"
	if req.Kind == "framework" {
		primary = "profile.fpl"
		subdir = "frameworks"
	}
	if req.Kind == "provider" {
		return "", fmt.Errorf("provider publish requires manifest.json via GitOps")
	}
	dir := filepath.Join(s.baseDir, "artifacts", subdir, req.Name, req.Version)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	fplPath := filepath.Join(dir, primary)
	if err := os.WriteFile(fplPath, []byte(req.PolicyFPL), 0o644); err != nil {
		return "", err
	}
	if req.PolicyYAML != "" {
		_ = os.WriteFile(filepath.Join(dir, "policy.yaml"), []byte(req.PolicyYAML), 0o644)
	}
	if req.PolicyJSON != "" {
		_ = os.WriteFile(filepath.Join(dir, "policy.json"), []byte(req.PolicyJSON), 0o644)
	}
	if req.ReadmeMarkdown != "" {
		_ = os.WriteFile(filepath.Join(dir, "README.md"), []byte(req.ReadmeMarkdown), 0o644)
	}
	rel := filepath.ToSlash(filepath.Join("artifacts", subdir, req.Name, req.Version, primary))
	return rel, nil
}

func (s *Server) routePolicies(w http.ResponseWriter, r *http.Request) {
	s.routeArtifact(w, r, "policy", s.index.Policies, "policy.fpl")
}

func (s *Server) routeFrameworks(w http.ResponseWriter, r *http.Request) {
	s.routeArtifact(w, r, "framework", s.index.Frameworks, "profile.fpl")
}

func (s *Server) routeProviders(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/v1/providers/")
	if rest == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if strings.HasSuffix(rest, "/versions") && !strings.Contains(rest, "/versions/") {
		name := strings.TrimSuffix(rest, "/versions")
		name = strings.TrimSuffix(name, "/")
		s.serveProviderVersionList(w, name)
		return
	}
	if name, version, ok := splitNameVersion(rest); ok {
		r2 := *r
		r2.SetPathValue("name", name)
		r2.SetPathValue("version", version)
		s.handleProviderVersion(w, &r2)
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
}

func (s *Server) routeArtifact(w http.ResponseWriter, r *http.Request, kind string, entries []catalog.Entry, primary string) {
	prefix := "/v1/policies/"
	if kind == "framework" {
		prefix = "/v1/frameworks/"
	}
	rest := strings.TrimPrefix(r.URL.Path, prefix)
	if rest == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if strings.HasSuffix(rest, "/versions") && !strings.Contains(rest, "/versions/") {
		name := strings.TrimSuffix(rest, "/versions")
		name = strings.TrimSuffix(name, "/")
		s.serveVersionList(w, r, kind, entries, name)
		return
	}
	if name, version, ok := splitNameVersion(rest); ok {
		for _, ent := range entries {
			if ent.Name != name {
				continue
			}
			path, ok := ent.Versions[version]
			if !ok {
				break
			}
			s.serveFPL(w, r, kind, name, version, path, ent.Description, ent.TrustTier, primary)
			return
		}
		http.Error(w, kind+" not found", http.StatusNotFound)
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
}

func splitNameVersion(rest string) (name, version string, ok bool) {
	i := strings.Index(rest, "/versions/")
	if i < 0 {
		return "", "", false
	}
	return rest[:i], rest[i+len("/versions/"):], true
}

func (s *Server) serveProviderVersionList(w http.ResponseWriter, name string) {
	for _, ent := range s.index.Providers {
		if ent.Name != name {
			continue
		}
		var vers []api.VersionEntry
		for v := range ent.Versions {
			vers = append(vers, api.VersionEntry{Version: v})
		}
		writeJSON(w, api.VersionsResponse{APIVersion: APIVersion, Kind: "provider", Name: name, Versions: vers})
		return
	}
	http.Error(w, "provider not found", http.StatusNotFound)
}

func (s *Server) handleLegacyPack(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.PathValue("name"))
	version := strings.TrimSpace(r.PathValue("version"))
	for _, p := range s.index.Packs {
		if p.Name != name {
			continue
		}
		path, ok := p.Versions[version]
		if !ok {
			break
		}
		s.serveFPL(w, r, "policy", name, version, path, p.Description, p.TrustTier, "policy.fpl")
		return
	}
	http.Error(w, "pack not found", http.StatusNotFound)
}

func (s *Server) serveVersionList(w http.ResponseWriter, r *http.Request, kind string, entries []catalog.Entry, name string) {
	name = strings.TrimSpace(name)
	for _, ent := range entries {
		if ent.Name != name {
			continue
		}
		var vers []api.VersionEntry
		for v, p := range ent.Versions {
			ch := ""
			if art, err := catalog.LoadFPLArtifact(p, fplPrimary(kind)); err == nil {
				ch = art.Meta.Changelog
			}
			vers = append(vers, api.VersionEntry{Version: v, Changelog: ch})
		}
		writeJSON(w, api.VersionsResponse{APIVersion: APIVersion, Kind: kind, Name: name, Versions: vers})
		return
	}
	http.Error(w, kind+" not found", http.StatusNotFound)
}

func fplPrimary(kind string) string {
	if kind == "framework" {
		return "profile.fpl"
	}
	return "policy.fpl"
}

func (s *Server) serveFPL(w http.ResponseWriter, r *http.Request, kind, name, version, path, desc, tier, primary string) {
	art, err := catalog.LoadFPLArtifact(path, primary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if strings.TrimSpace(desc) == "" && art.Meta.Description != "" {
		desc = art.Meta.Description
	}
	resp := api.PackVersionResponse{
		APIVersion: APIVersion, Kind: kind, Name: name, Version: version,
		Description: desc, PolicyFPL: string(art.FPLBytes),
		SHA256Hex: art.SHA256Hex, TrustTier: tier,
		Changelog: art.Meta.Changelog, ReadmeMarkdown: string(art.Readme),
		FarameshVersion: art.Meta.FarameshVersion, Compatibility: art.Meta.Compatibility,
		ParameterSchema: art.Meta.ParameterSchema,
	}
	if art.Meta.RulesSummary != nil {
		resp.RulesSummary = &api.RulesSummary{
			Permit: art.Meta.RulesSummary.Permit,
			Deny:   art.Meta.RulesSummary.Deny,
			Defer:  art.Meta.RulesSummary.Defer,
		}
	}
	if len(art.YAMLBytes) > 0 {
		resp.PolicyYAML = string(art.YAMLBytes)
	}
	if len(art.JSONBytes) > 0 {
		resp.PolicyJSON = string(art.JSONBytes)
	}
	if art.Signature != nil {
		resp.Signature = &api.ArtifactSignature{
			Algorithm: art.Signature.Algorithm, KeyID: art.Signature.KeyID,
			ValueB64: art.Signature.ValueB64, PublicKeyB64: art.Signature.PublicKeyB64,
		}
	}
	format := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("format")))
	switch format {
	case "yaml":
		if resp.PolicyYAML == "" {
			http.Error(w, "no YAML sidecar for this version", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/yaml")
		_, _ = w.Write([]byte(resp.PolicyYAML))
		return
	case "json":
		if resp.PolicyJSON == "" {
			http.Error(w, "no JSON sidecar for this version", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(resp.PolicyJSON))
		return
	case "fpl", "":
		writeJSON(w, resp)
	default:
		http.Error(w, "format must be fpl, yaml, or json", http.StatusBadRequest)
	}
}

func (s *Server) handleProviderVersion(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.PathValue("name"))
	version := strings.TrimSpace(r.PathValue("version"))
	for _, ent := range s.index.Providers {
		if ent.Name != name {
			continue
		}
		manifestPath, ok := ent.Versions[version]
		if !ok {
			break
		}
		man, raw, err := catalog.LoadProviderManifest(manifestPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		host := requestBaseURL(r)
		downloads := make(map[string]api.ProviderDownload, len(man.Downloads))
		for k, dl := range man.Downloads {
			url := dl.URL
			if strings.HasPrefix(url, "file://") {
				rel := strings.TrimPrefix(url, "file://")
				if !filepath.IsAbs(rel) {
					rel = filepath.Join(filepath.Dir(manifestPath), rel)
				}
				artDir := filepath.Join(s.baseDir, "artifacts")
				relFromArt, _ := filepath.Rel(artDir, rel)
				url = host + "/artifacts/" + filepath.ToSlash(relFromArt)
			}
			downloads[k] = api.ProviderDownload{URL: url, SHA256Hex: dl.SHA256Hex, Size: dl.Size}
		}
		readme, _ := os.ReadFile(filepath.Join(filepath.Dir(manifestPath), "README.md"))
		resp := api.ProviderVersionResponse{
			APIVersion: APIVersion, Kind: "provider", Name: name, Version: version,
			TrustTier: ent.TrustTier, Capabilities: ent.Capabilities,
			Downloads: downloads, ReadmeMarkdown: string(readme),
			ConfigSchema: man.ConfigSchema,
		}
		if man.Signature != nil {
			resp.Signature = &api.ArtifactSignature{
				Algorithm: man.Signature.Algorithm, KeyID: man.Signature.KeyID,
				ValueB64: man.Signature.ValueB64, PublicKeyB64: man.Signature.PublicKeyB64,
			}
		}
		resp.DevOnly = man.DevOnly
		_ = raw
		writeJSON(w, resp)
		return
	}
	http.Error(w, "provider not found", http.StatusNotFound)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

const APIVersion = api.APIVersion

// DefaultPlatformKey is the download map key for the current OS/arch.
var DefaultPlatformKey = runtime.GOOS + "_" + runtime.GOARCH
