import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    domains: ['lh3.googleusercontent.com'], // Googleプロフィール画像用
  },
  // Enable standalone output for Docker production builds
  output: 'standalone',
  // Development server configuration for Docker
  ...(process.env.NODE_ENV === 'development' && {
    experimental: {
      serverComponentsExternalPackages: []
    }
  })
};

export default nextConfig;
