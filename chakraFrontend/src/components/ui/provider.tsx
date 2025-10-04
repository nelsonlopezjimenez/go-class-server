'use client';

import { ChakraProvider } from '@chakra-ui/react';

import uiSystem from '@/lib/styles/theme/theme';

import { ColorModeProvider } from './color-mode';

export function Provider(props: React.PropsWithChildren) {
  return (
    <ChakraProvider value={uiSystem}>
      <ColorModeProvider>{props.children}</ColorModeProvider>
    </ChakraProvider>
  );
}
