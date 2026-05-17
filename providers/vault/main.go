package main

import (
	"fmt"
	"os"

	"github.com/faramesh/faramesh-registry/providers/sidecar"
)

func main() {
	if err := sidecar.Serve(newVaultServer()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
