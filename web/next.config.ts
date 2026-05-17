import type { NextConfig } from "next";
import path from "path";

const registryAPI =
  process.env.REGISTRY_API_URL || "http://127.0.0.1:9876";

const nextConfig: NextConfig = {
  output: "standalone",
  outputFileTracingRoot: path.join(__dirname),
  env: {
    REGISTRY_API_URL: process.env.REGISTRY_API_URL || "http://127.0.0.1:9876",
  },
  async rewrites() {
    return [
      {
        source: "/api/registry/:path*",
        destination: `${registryAPI}/:path*`,
      },
    ];
  },
};

export default nextConfig;
