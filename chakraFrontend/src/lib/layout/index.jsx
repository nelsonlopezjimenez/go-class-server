import { Box } from '@chakra-ui/react';

import Footer from './components/Footer';
import Nav from './components/Navigation';
import { Meta } from './components/meta';

import '@/lib/styles/stylesheets/global.css';

export const Layout = ({ children }) => {
  return (
    <Box
      maxWidth={'screen'}
      overflow="hidden"
      height="100vh"
      transition="5ms ease-out"
    >
      <Meta />
      <Nav />
      <Box
        // overflow='scroll'
        width="full"
        paddingX="2em"
        style={{
          maxHeight: 'calc(100vh - 3em)',
          overflow: 'auto',
          paddingTop: '1em',
        }}
        as="main"
      >
        {children}
        <Footer />
      </Box>
    </Box>
  );
};
