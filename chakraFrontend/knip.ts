import type { KnipConfig } from 'knip';

const config: KnipConfig = {
  project: ['src/**/*.{ts,tsx,js,jsx,css,scss}'],
  ignore: ['src/components/ui/*', 'src/lib/styles/theme/*','src\api\hooks.ts'],
  ignoreBinaries: ['changelogithub', 'tar', 'typegen'],
  ignoreDependencies: [
    '@react-router/node',
    'isbot',
    '@emotion/styled',
    'lodash',
    'next-themes',
  ],
  'react-router': true,
  rollup: true,
};

export default config;
