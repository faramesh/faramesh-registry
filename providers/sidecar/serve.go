// Package sidecar runs a ProviderService gRPC server on a Unix socket.
package sidecar

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	providerv1 "github.com/faramesh/faramesh-core/proto/provider/v1"
	"google.golang.org/grpc"
)

// Serve listens on socketPath (or FARAMESH_PROVIDER_SOCKET) and serves svc until SIGTERM.
func Serve(svc providerv1.ProviderServiceServer) error {
	socketPath := resolveSocketPath()
	if socketPath == "" {
		return fmt.Errorf("unix socket path required (argv[1] or FARAMESH_PROVIDER_SOCKET)")
	}
	_ = os.Remove(socketPath)
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("listen %s: %w", socketPath, err)
	}
	srv := grpc.NewServer()
	providerv1.RegisterProviderServiceServer(srv, svc)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(lis)
	}()
	select {
	case <-ctx.Done():
		srv.GracefulStop()
		_ = os.Remove(socketPath)
		return nil
	case err := <-errCh:
		return err
	}
}

func resolveSocketPath() string {
	if len(os.Args) > 1 && strings.TrimSpace(os.Args[1]) != "" {
		return strings.TrimSpace(os.Args[1])
	}
	return strings.TrimSpace(os.Getenv("FARAMESH_PROVIDER_SOCKET"))
}
