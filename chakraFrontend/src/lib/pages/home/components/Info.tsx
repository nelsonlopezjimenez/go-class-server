import { Card } from '@chakra-ui/react/card';
import { GridItem } from '@chakra-ui/react/grid';
import { Text } from '@chakra-ui/react/typography';
// import _ from 'lodash';
// import { useNavigate } from 'react-router';

const Info = ({ item }: { item: any }) => {
  // const navigate = useNavigate();
  // let path = '';
  // if (item.isQuarter) {
  //   path = `section/${_.kebabCase(item.title)}`;
  // } else {
  //   path = item.title.toLowerCase();
  // }

  return (
    <GridItem>
      <Card.Root
        colorPalette={{ _dark: 'blue', _light: 'black' }}
        // onClick={() => navigate(path)}
        size="sm"
        // variant={{
        //   base: 'elevated',
        //   _hover: 'subtle',
        // }}
        // _hover={{
        //   _dark: { backgroundColor: 'gray.emphasized/40' },
        //   _light: { backgroundColor: 'gray.emphasized/80' },
        // }}
        height="full"
      >
        <Card.Header>
          <Card.Title>{item.title}</Card.Title>
        </Card.Header>
        <Card.Body>
          <Text>{item.description}</Text>
        </Card.Body>
        <Card.Footer>
          {/* <Em color='colorPalette.solid'>
            Click to view this quarter's syllabus
          </Em> */}
        </Card.Footer>
      </Card.Root>
    </GridItem>
  );
};

export default Info;
