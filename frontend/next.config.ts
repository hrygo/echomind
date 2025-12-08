import type { NextConfig } from "next";

const isDesktop = process.env.BUILD_MODE === 'desktop';

const nextConfig: NextConfig = {
  output: isDesktop ? "export" : "standalone",
  images: {
    unoptimized: isDesktop,
  },
  async rewrites() {
    if (isDesktop) return [];
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ];
  },
  async redirects() {
    if (isDesktop) return [];
    return [
      {
        source: '/login',
        destination: '/auth',
        permanent: false,
      },
      {
        source: '/register',
        destination: '/auth?mode=register',
        permanent: false,
      },
    ];
  },
};

export default nextConfig;
