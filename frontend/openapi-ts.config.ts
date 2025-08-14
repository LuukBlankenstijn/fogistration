import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: '../go/docs/swagger.json',
  output: './src/clients/generated-client',
  plugins: [
    '@hey-api/client-fetch',
    '@hey-api/schemas',
    '@hey-api/sdk',
    { name: '@hey-api/transformers', dates: true },
    {
      name: '@hey-api/sdk',
      transformer: true,
    },
    '@tanstack/react-query',
    {
      enums: 'typescript',
      name: '@hey-api/typescript',
    },
  ],
});
