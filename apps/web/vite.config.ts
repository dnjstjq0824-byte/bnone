import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import { execSync } from 'child_process'
import fs from 'fs'

let commitHash = 'unknown'
try {
  commitHash = execSync('git rev-parse --short HEAD').toString().trim()
} catch (e) {}

const buildTime = new Date().toISOString()

// Ensure public/version.json is written at build time
try {
  const publicDir = path.resolve(__dirname, './public')
  if (!fs.existsSync(publicDir)) {
    fs.mkdirSync(publicDir, { recursive: true })
  }
  fs.writeFileSync(
    path.join(publicDir, 'version.json'),
    JSON.stringify(
      {
        app: 'MOA',
        version: '1.0.0',
        commit: commitHash,
        built_at: buildTime,
        mode: process.env.VITE_MOCK_ENABLED === 'true' ? 'static-mock-demo' : 'live-api',
      },
      null,
      2
    )
  )
} catch (e) {}

export default defineConfig({
  plugins: [vue()],
  define: {
    __COMMIT_HASH__: JSON.stringify(commitHash),
    __BUILD_TIME__: JSON.stringify(buildTime),
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
    },
  },
})
