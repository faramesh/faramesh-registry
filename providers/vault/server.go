package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	providerv1 "github.com/faramesh/faramesh-core/proto/provider/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type vaultConfig struct {
	addr      string
	token     string
	mount     string
	namespace string
	timeout   time.Duration
}

type vaultServer struct {
	providerv1.UnimplementedProviderServiceServer
	cfg    vaultConfig
	client *http.Client
}

func newVaultServer() *vaultServer {
	return &vaultServer{}
}

func (s *vaultServer) Init(ctx context.Context, req *providerv1.InitRequest) (*providerv1.ProviderInfo, error) {
	cfg := req.GetConfig()
	addr := strings.TrimSpace(cfg["addr"])
	token := strings.TrimSpace(cfg["token"])
	if addr == "" || token == "" {
		return nil, fmt.Errorf("addr and token are required")
	}
	mount := strings.TrimSpace(cfg["mount"])
	if mount == "" {
		mount = "secret"
	}
	s.cfg = vaultConfig{
		addr:      strings.TrimRight(addr, "/"),
		token:     token,
		mount:     mount,
		namespace: strings.TrimSpace(cfg["namespace"]),
		timeout:   10 * time.Second,
	}
	s.client = &http.Client{Timeout: s.cfg.timeout}
	if req.GetDryRun() {
		return vaultInfo(true, "dry-run ok"), nil
	}
	if err := s.pingHealth(ctx); err != nil {
		return nil, err
	}
	return vaultInfo(true, "connected"), nil
}

func vaultInfo(healthy bool, detail string) *providerv1.ProviderInfo {
	return &providerv1.ProviderInfo{
		Capabilities: []providerv1.Capability{providerv1.Capability_CAPABILITY_SECRETS},
		Health:       &providerv1.HealthStatus{Healthy: healthy, Detail: detail},
		Version:      "1.0.0",
		Schema: []*providerv1.ConfigSchemaField{
			{Name: "addr", Type: "string", Required: true},
			{Name: "token", Type: "string", Required: true},
			{Name: "mount", Type: "string", Required: false},
			{Name: "namespace", Type: "string", Required: false},
		},
	}
}

func (s *vaultServer) HealthCheck(ctx context.Context, _ *providerv1.HealthRequest) (*providerv1.HealthStatus, error) {
	if s.client == nil {
		return &providerv1.HealthStatus{Healthy: false, Detail: "not initialized"}, nil
	}
	if err := s.pingHealth(ctx); err != nil {
		return &providerv1.HealthStatus{Healthy: false, Detail: err.Error()}, nil
	}
	return &providerv1.HealthStatus{Healthy: true, Detail: "ok"}, nil
}

func (s *vaultServer) GetSecret(ctx context.Context, req *providerv1.SecretRequest) (*providerv1.Secret, error) {
	if s.client == nil {
		return nil, status.Error(codes.FailedPrecondition, "not initialized")
	}
	path := strings.TrimSpace(req.GetPath())
	if path == "" {
		return nil, status.Error(codes.InvalidArgument, "path is required")
	}
	vaultPath := resolveVaultPath(s.cfg.mount, path)
	url := fmt.Sprintf("%s/v1/%s", s.cfg.addr, vaultPath)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "build request: %v", err)
	}
	httpReq.Header.Set("X-Vault-Token", s.cfg.token)
	if s.cfg.namespace != "" {
		httpReq.Header.Set("X-Vault-Namespace", s.cfg.namespace)
	}
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "vault request: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, status.Errorf(codes.Internal, "vault %s: %d %s", vaultPath, resp.StatusCode, truncate(string(body), 200))
	}
	var vaultResp vaultSecretResponse
	if err := json.Unmarshal(body, &vaultResp); err != nil {
		return nil, status.Errorf(codes.Internal, "parse response: %v", err)
	}
	value, leaseID := extractCredentialValue(vaultResp)
	if value == "" {
		return nil, status.Errorf(codes.Internal, "no credential value at %s", vaultPath)
	}
	out := &providerv1.Secret{Value: []byte(value), Version: "vault"}
	if vaultResp.LeaseDuration > 0 {
		out.Ttl = durationpb.New(time.Duration(vaultResp.LeaseDuration) * time.Second)
	}
	_ = leaseID
	return out, nil
}

func resolveVaultPath(mount, toolPath string) string {
	switch mount {
	case "aws", "database", "pki":
		role := strings.ReplaceAll(toolPath, "/", "-")
		return fmt.Sprintf("%s/creds/%s", mount, role)
	default:
		return fmt.Sprintf("%s/data/faramesh/%s", mount, toolPath)
	}
}

type vaultSecretResponse struct {
	Data          map[string]any `json:"data"`
	LeaseID       string         `json:"lease_id"`
	LeaseDuration int            `json:"lease_duration"`
}

func extractCredentialValue(resp vaultSecretResponse) (value string, leaseID string) {
	leaseID = resp.LeaseID
	if resp.Data == nil {
		return "", leaseID
	}
	if inner, ok := resp.Data["data"].(map[string]any); ok {
		for _, key := range []string{"value", "api_key", "token", "password"} {
			if v, ok := inner[key].(string); ok && v != "" {
				return v, leaseID
			}
		}
		for _, v := range inner {
			if s, ok := v.(string); ok && s != "" {
				return s, leaseID
			}
		}
	}
	for _, key := range []string{"value", "api_key", "token", "password", "access_key", "secret_key"} {
		if v, ok := resp.Data[key].(string); ok && v != "" {
			return v, leaseID
		}
	}
	return "", leaseID
}

func (s *vaultServer) pingHealth(ctx context.Context) error {
	url := s.cfg.addr + "/v1/sys/health"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Set("X-Vault-Token", s.cfg.token)
	if s.cfg.namespace != "" {
		httpReq.Header.Set("X-Vault-Namespace", s.cfg.namespace)
	}
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("vault health: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("vault unhealthy: %d %s", resp.StatusCode, truncate(string(body), 120))
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func (s *vaultServer) VerifyIdentity(context.Context, *providerv1.Identity) (*providerv1.VerificationResult, error) {
	return nil, status.Error(codes.Unimplemented, "IDENTITY not supported")
}

func (s *vaultServer) Sign(context.Context, *providerv1.SignRequest) (*providerv1.Signature, error) {
	return nil, status.Error(codes.Unimplemented, "KMS not supported")
}

func (s *vaultServer) SinkDPR(context.Context, *providerv1.DPRRecord) (*providerv1.SinkAck, error) {
	return nil, status.Error(codes.Unimplemented, "AUDIT_SINK not supported")
}

func (s *vaultServer) CostEstimate(context.Context, *providerv1.CostRequest) (*providerv1.CostEstimateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "COST not supported")
}
