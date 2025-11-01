import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  output: 'standalone',
  /*
  eslint: {
    // Warning: This allows production builds to successfully complete even if
    // your project has ESLint errors.
    ignoreDuringBuilds: true,
  },
  */
  async rewrites() {
    // This proxy is used ONLY in local development.
    // In production, the frontend will call the backend API directly using the NEXT_PUBLIC_API_URL.
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*', // Proxy to the local Go backend
      },
    ];
  },
};

export default nextConfig;
