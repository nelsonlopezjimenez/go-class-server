import { Box, Button, Grid, Image } from '@chakra-ui/react';
import { useNavigate } from 'react-router';

const Page404 = () => {
  const navigate = useNavigate();

  const handleBackToHome = () => navigate('/');

  return (
    <Grid gap={4} textAlign="center">
      <Box maxWidth={[280, 400]} marginX="auto">
        <Image width={400} src="/assets/404 Error-rafiki.svg" />
      </Box>

      <Box>
        <Button onClick={handleBackToHome}>Let&apos;s Head Back</Button>
      </Box>
    </Grid>
  );
};

export default Page404;
