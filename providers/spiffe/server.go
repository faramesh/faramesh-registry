package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	providerv1 "github.com/faramesh/faramesh-core/proto/provider/v1"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type spiffeServer struct {
	providerv1.UnimplementedProviderServiceServer
	socket      string
	trustDomain string
}

func newSPIFFEServer() *spiffeServer {
	return &spiffeServer{}
}

func (s *spiffeServer) Init(_ context.Context, req *providerv1.InitRequest) (*providerv1.ProviderInfo, error) {
	cfg := req.GetConfig()
	s.socket = strings.TrimSpace(cfg["socket"])
	s.trustDomain = strings.TrimSpace(cfg["trust_domain"])
	if s.socket == "" || s.trustDomain == "" {
		return nil, fmt.Errorf("socket and trust_domain are required")
	}
	if req.GetDryRun() {
		return spiffeInfo(true, "dry-run ok"), nil
	}
	if err := s.checkSocket(context.Background()); err != nil {
		return nil, err
	}
	return spiffeInfo(true, "initialized"), nil
}

func spiffeInfo(healthy bool, detail string) *providerv1.ProviderInfo {
	return &providerv1.ProviderInfo{
		Capabilities: []providerv1.Capability{providerv1.Capability_CAPABILITY_IDENTITY},
		Health:       &providerv1.HealthStatus{Healthy: healthy, Detail: detail},
		Version:      "1.0.0",
		Schema: []*providerv1.ConfigSchemaField{
			{Name: "socket", Type: "string", Required: true},
			{Name: "trust_domain", Type: "string", Required: true},
		},
	}
}

func (s *spiffeServer) HealthCheck(ctx context.Context, _ *providerv1.HealthRequest) (*providerv1.HealthStatus, error) {
	if s.socket == "" {
		return &providerv1.HealthStatus{Healthy: false, Detail: "not initialized"}, nil
	}
	if err := s.checkSocket(ctx); err != nil {
		return &providerv1.HealthStatus{Healthy: false, Detail: err.Error()}, nil
	}
	return &providerv1.HealthStatus{Healthy: true, Detail: "ok"}, nil
}

func (s *spiffeServer) VerifyIdentity(ctx context.Context, req *providerv1.Identity) (*providerv1.VerificationResult, error) {
	if s.socket == "" {
		return nil, status.Error(codes.FailedPrecondition, "not initialized")
	}
	client, err := workloadapi.New(ctx, workloadapi.WithAddr(normalizeSPIFFESocket(s.socket)))
	if err != nil {
		return &providerv1.VerificationResult{Valid: false, Reason: fmt.Sprintf("connect: %v", err)}, nil
	}
	defer client.Close()
	svid, err := client.FetchX509SVID(ctx)
	if err != nil {
		return &providerv1.VerificationResult{Valid: false, Reason: fmt.Sprintf("fetch svid: %v", err)}, nil
	}
	if svid == nil {
		return &providerv1.VerificationResult{Valid: false, Reason: "nil svid"}, nil
	}
	id := svid.ID.String()
	if !strings.HasPrefix(strings.ToLower(id), "spiffe://") {
		return &providerv1.VerificationResult{Valid: false, Reason: "invalid spiffe id"}, nil
	}
	td := trustDomainFromSPIFFEID(id)
	if !strings.EqualFold(td, s.trustDomain) {
		return &providerv1.VerificationResult{
			Valid:  false,
			Reason: fmt.Sprintf("trust domain %q does not match config %q", td, s.trustDomain),
		}, nil
	}
	if want := strings.TrimSpace(req.GetId()); want != "" && !strings.EqualFold(want, id) {
		return &providerv1.VerificationResult{Valid: false, Reason: "identity mismatch"}, nil
	}
	var expires *timestamppb.Timestamp
	if len(svid.Certificates) > 0 {
		expires = timestamppb.New(svid.Certificates[0].NotAfter)
	}
	return &providerv1.VerificationResult{
		Valid:     true,
		Subject:   id,
		ExpiresAt: expires,
	}, nil
}

func (s *spiffeServer) checkSocket(ctx context.Context) error {
	client, err := workloadapi.New(ctx, workloadapi.WithAddr(normalizeSPIFFESocket(s.socket)))
	if err != nil {
		return fmt.Errorf("spiffe socket: %w", err)
	}
	client.Close()
	return nil
}

func normalizeSPIFFESocket(socketPath string) string {
	socketPath = strings.TrimSpace(socketPath)
	if socketPath == "" {
		return socketPath
	}
	if strings.Contains(socketPath, "://") {
		return socketPath
	}
	if strings.HasPrefix(socketPath, "unix:") {
		return "unix://" + strings.TrimPrefix(socketPath, "unix:")
	}
	if strings.HasPrefix(socketPath, "/") {
		return "unix://" + filepath.Clean(socketPath)
	}
	return socketPath
}

func trustDomainFromSPIFFEID(id string) string {
	const prefix = "spiffe://"
	if !strings.HasPrefix(strings.ToLower(id), prefix) {
		return ""
	}
	rest := strings.TrimPrefix(id, prefix)
	if i := strings.Index(rest, "/"); i >= 0 {
		return rest[:i]
	}
	return rest
}

func (s *spiffeServer) GetSecret(context.Context, *providerv1.SecretRequest) (*providerv1.Secret, error) {
	return nil, status.Error(codes.Unimplemented, "SECRETS not supported")
}

func (s *spiffeServer) Sign(context.Context, *providerv1.SignRequest) (*providerv1.Signature, error) {
	return nil, status.Error(codes.Unimplemented, "KMS not supported")
}

func (s *spiffeServer) SinkDPR(context.Context, *providerv1.DPRRecord) (*providerv1.SinkAck, error) {
	return nil, status.Error(codes.Unimplemented, "AUDIT_SINK not supported")
}

func (s *spiffeServer) CostEstimate(context.Context, *providerv1.CostRequest) (*providerv1.CostEstimateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "COST not supported")
}
