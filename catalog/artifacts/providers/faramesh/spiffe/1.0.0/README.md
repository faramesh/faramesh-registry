# faramesh/spiffe

SPIRE workload identity `ProviderService` sidecar (IDENTITY capability).

```fpl
import "registry.faramesh.dev/providers/faramesh/spiffe@1.0.0"

provider "spiffe" {
  type         = "spiffe"
  socket       = "/run/spire/sockets/agent.sock"
  trust_domain = "example.org"
}
```
