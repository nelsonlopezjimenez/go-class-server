import { useGetInfoQuery } from '@/api/apiSlice';
import { Card } from '@chakra-ui/react/card';

export type resource = {
  name: string;
  url: string;
  description: string;
};

const NetworkResources: React.FC = () => {
  const { data, isSuccess } = useGetInfoQuery('lan-links');

  return isSuccess ? (
    <>
      <h1>Links to in class Resources</h1>

      <Card.Root mx="10px">
        <Card.Body
          id="localLinks"
          style={{
            // backgroundColor: "#99bbd9",
            // color: "#06548a",
            margin: '0 10%',
          }}
          className="markdown-content"
          dangerouslySetInnerHTML={{ __html: data }}
        />
      </Card.Root>
    </>
  ) : (
    <></>
  );
};

export default NetworkResources;
