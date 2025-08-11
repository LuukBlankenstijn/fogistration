import { defineConfig } from 'vite'
import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

import tanstackRouter from '@tanstack/router-plugin/vite'
import { fileURLToPath } from 'node:url'
import { dirname, resolve } from 'node:path'


const __dirname = dirname(fileURLToPath(import.meta.url))
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    tanstackRouter({ autoCodeSplitting: true }),
    viteReact(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080"
      }
    }
  }
})
