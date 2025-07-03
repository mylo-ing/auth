import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
    output: 'standalone',
    async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://host.docker.internal:3517/:path*',
      },
    ];
  },
};

export default nextConfig;
