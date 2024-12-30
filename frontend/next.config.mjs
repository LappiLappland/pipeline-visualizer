/** @type {import('next').NextConfig} */
const nextConfig = {
    experimental: {
        serverActions: {
            allowedOrigins: [
                'http://localhost:8080',
            ]
        },
    },
    output: 'standalone',
};

export default nextConfig;


