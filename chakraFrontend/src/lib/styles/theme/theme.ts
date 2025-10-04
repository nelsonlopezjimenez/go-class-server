import { createSystem, defaultConfig, defineConfig } from '@chakra-ui/react';

const uiConfig = defineConfig({
  strictTokens: true,
  cssVarsRoot: ':where(:root, :host)',
  cssVarsPrefix: 'ck',
  globalCss: {
    'html, body': {
      fontFamily: 'body',
      boxSizing: 'border-box',
    },
    /* width */
    '::-webkit-scrollbar': {
      width: '4',
    },

    /* Track */
    '::-webkit-scrollbar-track': {
      background: 'currentBg',
      // height: '80%',
    },

    /* Handle */
    '::-webkit-scrollbar-thumb': {
      background: 'blue.muted',
    },

    /* Handle on hover */
    '::-webkit-scrollbar-thumb:hover': {
      background: 'blue.focusRing',
    },

    '.markdown-content h1': {
      color: 'blue.800',
    },
    '.markdown-content h2': {
      color: 'blue.700',
    },
    '.markdown-content h3': {
      color: 'blue.600',
    },
    '.markdown-content h4': {
      color: 'blue.500',
    },
    '.markdown-content p': {
      color: 'colorPalette.solid',
    },
    '.markdown-content code': {
      color: 'gray.800',
    },
    '.markdown-content pre': {
      color: 'gray.100',
    },
    '.markdown-content pre code': {
      color: 'gray.100',
    },
  },
  theme: {
    breakpoints: {
      sm: '320px',
      md: '768px',
      lg: '960px',
      xl: '1200px',
    },
    keyframes: {
      spin: {
        from: { transform: 'rotate(0deg)' },
        to: { transform: 'rotate(360deg)' },
      },
    },
    tokens: {
      colors: {
        primary: { value: '#0FEE0F' },
        secondary: { value: '#EE0F0F' },
        red: { value: '#EE0F0F' },
      },
      gradients: {
        primary: {
          value: 'linear-gradient(to bottom right, #222,rgb(34, 47, 62))',
        },
      },
      fonts: {
        heading: { value: 'Plus Jakarta Sans Variable, sans-serif' },
        body: { value: 'Plus Jakarta Sans Variable, sans-serif' },
      },
    },
    semanticTokens: {
      colors: {
        danger: { value: '{colors.red}' },
      },
    },
    textStyles: {},
    layerStyles: {},
    animationStyles: {},
    recipes: {},
    slotRecipes: {},
  },
});

const uiSystem = createSystem(defaultConfig, uiConfig);

export default uiSystem;
