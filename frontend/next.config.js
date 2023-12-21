/** @type {import('next').NextConfig} */
const nextConfig = {
    basePath: process.env.NEXT_PUBLIC_BASEPATH,
    rewrites: async () => [
        {
            source: '/api/:path*',
            destination: 'http://localhost:9998/api/:path*',
        },
    ],
}

module.exports = nextConfig
