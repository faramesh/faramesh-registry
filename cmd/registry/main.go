// Command registry serves the official Faramesh Registry HTTP API (v1).
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/faramesh/faramesh-registry/internal/server"
)

func main() {
	catalogDir := flag.String("catalog", "catalog", "path to catalog directory (catalog.json + artifacts)")
	listen := flag.String("listen", "127.0.0.1:9876", "listen address")
	flag.Parse()

	srv, err := server.New(*catalogDir)
	if err != nil {
		log.Fatalf("registry: %v", err)
	}
	log.Printf("Faramesh Registry listening on http://%s", *listen)
	log.Printf("  well-known: http://%s/.well-known/faramesh.json", *listen)
	log.Fatal(http.ListenAndServe(*listen, srv.Handler()))
}
