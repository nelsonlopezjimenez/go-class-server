/* eslint-disable import/no-extraneous-dependencies */
import { reactRouter } from '@react-router/dev/vite';
import { defineConfig } from 'vite';
import checker from 'vite-plugin-checker';
import type { VitePWAOptions } from 'vite-plugin-pwa';
import { VitePWA } from 'vite-plugin-pwa';
import tsconfigPaths from 'vite-tsconfig-paths';

const pwaOptions: Partial<VitePWAOptions> = {
  registerType: 'autoUpdate',
  manifest: {
    short_name: 'class-server',
    name: 'CIS Class Server',
    lang: 'en',
    start_url: './',
    background_color: '#FFFFFF',
    theme_color: '#FFFFFF',
    dir: 'ltr',
    display: 'standalone',
    prefer_related_applications: false,
    icons: [
      {
        src: './assets/favicon.svg',
        purpose: 'any',
        sizes: '48x48 72x72 96x96 128x128 256x256',
      },
    ],
  },
};

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
