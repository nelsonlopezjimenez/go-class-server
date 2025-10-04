/* eslint-disable import/no-extraneous-dependencies */
import { reactRouter } from '@react-router/dev/vite';
import { defineConfig } from 'vite';
import checker from 'vite-plugin-checker';
import type { VitePWAOptions } from 'vite-plugin-pwa';
import { VitePWA } from 'vite-plugin-pwa';
import tsconfigPaths from 'vite-tsconfig-paths';

const pwaOptions: Partial<VitePWAOptions> = {
//   registerType: 'autoUpdate',
//   manifest: {
//     short_name: 'class-server',
//     name: 'CIS Class Server',
//     lang: 'en',
//     start_url: './',
//     background_color: '#FFFFFF',
//     theme_color: '#FFFFFF',
//     dir: 'ltr',
//     display: 'standalone',
//     prefer_related_applications: false,
//     icons: [
//       {
//         src: './assets/favicon.svg',
//         purpose: 'any',
//         sizes: '48x48 72x72 96x96 128x128 256x256',
//       },
//     ],
//   },
      registerType: "autoUpdate",
      strategies: "generateSW",
      workbox: {
        globDirectory: "build/client",
        globPatterns: [
          "**/*.{js,css,html,ico,png,svg,woff,woff2}",
        ],
        // Don't fail build if patterns don't match
        globIgnores: [
          "**/node_modules/**/*",
          "sw.js",
          "workbox-*.js"
        ],
        // Add runtime caching for better offline support
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/fonts\.googleapis\.com\/.*/i,
            handler: "CacheFirst",
            options: {
              cacheName: "google-fonts-cache",
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 60 * 60 * 24 * 365, // 1 year
              },
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
        ],
      },
      manifest: {
        name: "CIS Class Server",
        short_name: "CIS Server",
        description: "CIS Class Server Frontend",
        theme_color: "#ffffff",
        icons: [
          {
            src: "icon-192x192.png",
            sizes: "192x192",
            type: "image/png",
          },
          {
            src: "icon-512x512.png",
            sizes: "512x512",
            type: "image/png",
          },
        ],
      },
      // Don't fail build on warnings
      failOnWarning: false,
}


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    reactRouter(),
    checker({
      typescript: false,
      biome: true,
    }),
    tsconfigPaths(),
    VitePWA(pwaOptions),
  ],
  build: {
    //  Ensure sourcemaps are generated properly
    sourcemap: true,
    rollupOptions : {
      output: {
        //  Better chunk naming
        manualChunks: undefined,
      }
    }
  },
  preview: {
    port: 3002,

    proxy: {
      '/api': {
        target: 'http://127.0.0.1:22022',
      },
      '/websites': {
        target: 'http://127.0.0.1:22022',
      },
      '/lessons': {
        target: 'http://127.0.0.1:22022',
      },
    },
  },
  server: {
    port: 3002,

    proxy: {
      '/api': {
        target: 'http://127.0.0.1:22022',
      },
      '/raw': {
        target: 'http://127.0.0.1:22022',
      },
      '/websites': {
        target: 'http://127.0.0.1:22022',
      },
      '/lessons': {
        target: 'http://127.0.0.1:22022',
      },
      '/images': {
        target: 'http://127.0.0.1:22022',
      },
    },
  },
  ssr: {
    // noExternal: ['react-helmet-async'], // temporary
  },
});
