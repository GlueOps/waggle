import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
    input: './docs/openapi.json',
    output: {
        path: './ui/src/sdk',
    },
    plugins: [
        '@hey-api/client-fetch',
        '@hey-api/sdk',
        {
            name: '@tanstack/react-query',
            queryKeys: true,
            queryOptions: true,
            infiniteQueryKeys: true,
            infiniteQueryOptions: true,
            mutationOptions: true,
        },
    ],
});