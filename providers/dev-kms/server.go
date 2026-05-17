package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"sync"

	providerv1 "github.com/faramesh/faramesh-core/proto/provider/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type devKMSServer struct {
	providerv1.UnimplementedProviderServiceServer
	mu  sync.Mutex
	key *rsa.PrivateKey
}

func newDevKMSServer() *devKMSServer {
	return &devKMSServer{}
}

func (s *devKMSServer) Init(_ context.Context, req *providerv1.InitRequest) (*providerv1.ProviderInfo, error) {
	if req.GetDryRun() {
		return kmsInfo(), nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.key == nil {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		s.key = key
	}
	return kmsInfo(), nil
}

func kmsInfo() *providerv1.ProviderInfo {
	return &providerv1.ProviderInfo{
		Capabilities: []providerv1.Capability{providerv1.Capability_CAPABILITY_KMS},
		Health:       &providerv1.HealthStatus{Healthy: true, Detail: "ephemeral dev KMS"},
		Version:      "1.0.0",
	}
}

func (s *devKMSServer) HealthCheck(context.Context, *providerv1.HealthRequest) (*providerv1.HealthStatus, error) {
	return &providerv1.HealthStatus{Healthy: true, Detail: "ok"}, nil
}

func (s *devKMSServer) Sign(_ context.Context, req *providerv1.SignRequest) (*providerv1.Signature, error) {
	if req == nil || len(req.GetPayload()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "payload required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.key == nil {
		return nil, status.Error(codes.FailedPrecondition, "not initialized")
	}
	sum := sha256.Sum256(req.GetPayload())
	sig, err := rsa.SignPSS(rand.Reader, s.key, crypto.SHA256, sum[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign: %v", err)
	}
	return &providerv1.Signature{
		Signature: sig,
		Algorithm: "RSA-PSS-SHA256",
		KeyId:     "dev-kms",
	}, nil
}

func (s *devKMSServer) GetSecret(context.Context, *providerv1.SecretRequest) (*providerv1.Secret, error) {
	return nil, status.Error(codes.Unimplemented, "SECRETS not supported")
}

func (s *devKMSServer) VerifyIdentity(context.Context, *providerv1.Identity) (*providerv1.VerificationResult, error) {
	return nil, status.Error(codes.Unimplemented, "IDENTITY not supported")
}

func (s *devKMSServer) SinkDPR(context.Context, *providerv1.DPRRecord) (*providerv1.SinkAck, error) {
	return nil, status.Error(codes.Unimplemented, "AUDIT_SINK not supported")
}

func (s *devKMSServer) CostEstimate(context.Context, *providerv1.CostRequest) (*providerv1.CostEstimateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "COST not supported")
}
