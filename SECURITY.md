# Security Policy

The Faramesh registry distributes policy packs, framework profiles, and provider
artifacts consumed by Faramesh runtimes. Treat registry integrity issues as
security issues.

## Reporting a Vulnerability

Do not open a public issue for security vulnerabilities.

Preferred reporting paths:

- GitHub Security Advisory: https://github.com/faramesh/faramesh-registry/security/advisories/new
- Email: security@faramesh.dev

Include the affected artifact path, version, reproduction steps, impact, and any
suggested remediation.

## Response Timeline

- Acknowledgment: within 3 business days
- Initial assessment: within 7 business days
- Status updates: as remediation progresses

## Scope

Security reports may include:

- malicious or tampered registry artifacts;
- incorrect signatures, hashes, or provenance data;
- provider binary vulnerabilities;
- policy packs that create unexpected privilege expansion;
- registry server or resolver vulnerabilities.

## Supported Versions

The `main` branch and latest published catalog entries are supported. Older
catalog entries are handled on a best-effort basis.
