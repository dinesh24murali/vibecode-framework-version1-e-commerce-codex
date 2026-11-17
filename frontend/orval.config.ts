import { defineConfig } from 'orval'

export default defineConfig({
  api: {
    input: '../docs/02_outputs/04_api_spec.yaml',
    output: {
      target: './lib/api/generated',
      schemas: './lib/api/generated/model',
      client: 'fetch',
      mode: 'tags-split',
      clean: true,
      override: {
        mutator: { path: './lib/api/client.ts', name: 'apiFetch' },
      },
    },
  },
})
