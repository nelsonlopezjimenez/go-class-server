import { useGetLinksQuery } from '@/api/apiSlice';
import { DialogRoot, DialogTrigger } from '@/components/ui/dialog';
import { LinkButton } from '@/components/ui/link-button';
import { SimpleGrid } from '@chakra-ui/react';
import { Button } from '@chakra-ui/react/button';
import { Card } from '@chakra-ui/react/card';
import { Text } from '@chakra-ui/react/typography';
import WebsiteDownloadDialog from './components/WebsiteDownloadDialog';

type linkData = {
  Domain: string;
  IndexPath: string;
  Installed: boolean;
  IsCurrent: boolean;
};

export type ServerData = Array<linkData>;

const pathfinder = (link: linkData) => {
  const base = 'http://localhost:22022/';
  return `${base}${link.IndexPath ? link.IndexPath.slice(3) : link.Domain}`;
};

const OfflineLinks = () => {
  const { data, isSuccess } = useGetLinksQuery('links');
  return (
    <>
      {isSuccess && data.length > 0 ? (
        <>
          <Card.Root variant="elevated" mb="4">
            <Card.Header>
              <Card.Title>Installed Websites</Card.Title>
            </Card.Header>
            <Card.Body>
              <SimpleGrid
                gap={1}
                columns={{ base: 1, sm: 2, md: 4, lg: 5, xl: 6 }}
              >
                {data.map((link: linkData, i: number) => {
                  return (
                    link.Installed && (
                      <LinkButton
                        key={link.Domain + i}
                        href={pathfinder(link)}
                        variant="surface"
                        colorPalette="blue"
                      >
                        <Text>{link.Domain}</Text>
                      </LinkButton>
                    )
                  );
                })}
              </SimpleGrid>
            </Card.Body>
          </Card.Root>
          <Card.Root variant="elevated">
            <Card.Header>
              <Card.Title>
                Websites Available to Install -- You must be connected to the
                network
              </Card.Title>
            </Card.Header>
            <Card.Body>
              <SimpleGrid
                gap={1}
                columns={{ base: 1, sm: 1, md: 3, lg: 4, xl: 5 }}
              >
                {data.map((link: linkData, i: number) => {
                  return (
                    !link.Installed && (
                      <Card.Root key={link.Domain + i} size="sm">
                        <Card.Body>
                          <Text>{link.Domain}</Text>
                        </Card.Body>
                        <Card.Footer justifyContent={'flex-end'}>
                          <DialogRoot role="alertdialog">
                            <DialogTrigger asChild>
                              <Button variant="surface" colorPalette="green">
                                Download
                              </Button>
                            </DialogTrigger>
                            <WebsiteDownloadDialog domain={link.Domain} />
                          </DialogRoot>
                        </Card.Footer>
                      </Card.Root>
                    )
                  );
                })}
              </SimpleGrid>
            </Card.Body>
          </Card.Root>
        </>
      ) : (
        <p>Nothing here</p>
      )}
    </>
  );
};

export default OfflineLinks;
