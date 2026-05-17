.PHONY: gen-key providers sign-catalog sign-all validate

CATALOG ?= catalog
KEY_ID ?= faramesh-ed25519-2026
PROVIDERS := vault spiffe dev-kms
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

gen-key:
	go run ./cmd/gen-signing-key

providers:
	cd providers && go mod tidy
	@for p in $(PROVIDERS); do \
		echo "==> $$p"; \
		art=catalog/artifacts/providers/faramesh/$$p/1.0.0/bin; \
		mkdir -p $$art; \
		for plat in $(PLATFORMS); do \
			os=$${plat%/*}; arch=$${plat#*/}; key=$${os}_$$arch; \
			out=$$art/faramesh-provider-$$p-$$key; \
			( cd providers && GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" \
				-o ../$$out ./$$p ) && echo "  $$out"; \
		done; \
	done

sign-catalog:
	@test -n "$$REGISTRY_SIGNING_KEY_B64" || (echo "set REGISTRY_SIGNING_KEY_B64" >&2; exit 1)
	go run ./cmd/sign-catalog -catalog $(CATALOG) -key-id $(KEY_ID)
	go run ./cmd/sign-catalog -catalog $(CATALOG) -key-id $(KEY_ID) -providers

sign-all: providers sign-catalog

validate:
	./scripts/validate-catalog.sh
