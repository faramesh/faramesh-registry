# Faramesh Registry Web

Next.js browse UI for `registry.faramesh.dev`.

- **FPL** is the default tab on policy/framework detail pages.
- **YAML / JSON** tabs appear when publishers ship sidecars.

## Dev

```bash
# Terminal 1 — API
cd .. && go run ./cmd/registry -catalog catalog

# Terminal 2 — web
npm install
npm run dev
```

Open http://localhost:3001

`REGISTRY_API_URL` defaults to `http://127.0.0.1:9876` (proxied via `next.config.ts` rewrites).
