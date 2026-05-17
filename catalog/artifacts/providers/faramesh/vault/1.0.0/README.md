# faramesh/vault

Official HashiCorp Vault `ProviderService` sidecar (secrets capability).

```fpl
import "registry.faramesh.dev/providers/faramesh/vault@1.0.0"

provider "vault" {
  type  = "vault"
  addr  = env("VAULT_ADDR")
  token = env("VAULT_TOKEN")
  mount = "secret"
}
```
