# Production deployment

## Prerequisites

- Render account (render.com)
- Netlify account (netlify.com)
- Your REGISTRY_SIGNING_KEY_B64 value (from gen-signing-key)
- DNS access for registry.faramesh.dev and api.registry.faramesh.dev

## Deploy the API to Render

1. Go to render.com/dashboard
2. New → Blueprint
3. Connect your GitHub repository
4. Set repository root to: faramesh-registry/
5. Render detects render.yaml automatically
6. When prompted, set REGISTRY_SIGNING_KEY_B64 to your private key value
7. Deploy
8. Note the generated URL (e.g. faramesh-registry-api.onrender.com)

## Deploy the frontend to Netlify

1. Go to app.netlify.com
2. Add new site → Import from Git
3. Connect your GitHub repository
4. Set base directory to: faramesh-registry/web
5. Build command: npm run build (auto-detected from netlify.toml)
6. Publish directory: .next (auto-detected from netlify.toml)
7. Add environment variable:
   REGISTRY_API_URL = https://faramesh-registry-api.onrender.com
8. Deploy
9. Note the generated URL (e.g. faramesh-registry.netlify.app)

## Set custom domains

**Render (API):**
- Dashboard → faramesh-registry-api → Settings → Custom Domains
- Add: api.registry.faramesh.dev
- Add the CNAME Render shows you to your DNS provider

**Netlify (Frontend):**
- Dashboard → your site → Domain settings
- Add custom domain: registry.faramesh.dev
- Add the CNAME Netlify shows you to your DNS provider

## Verify

Once DNS propagates (5–60 minutes):

```bash
curl https://api.registry.faramesh.dev/.well-known/faramesh.json
curl https://api.registry.faramesh.dev/v1/search?q=stripe
```

Open https://registry.faramesh.dev in a browser and confirm the catalog loads.

## Update the faramesh CLI registry URL

Once live, engineers using faramesh set:

```bash
export FARAMESH_REGISTRY_URL=https://api.registry.faramesh.dev
```

Or in governance.fms runtime block:

```fpl
runtime {
  horizon {
    enabled      = true
    registry_url = "https://api.registry.faramesh.dev"
  }
}
```
