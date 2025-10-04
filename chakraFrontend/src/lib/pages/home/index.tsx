import { SimpleGrid } from '@chakra-ui/react';

import trident from '@/incomming/assets/EdmondsTrident.png';
import tridentLight from '@/incomming/assets/EdmondsTridentLight.png';

import { useColorMode } from '@/components/ui/color-mode';
import DATA from '@/incomming/data/courseInfo.js';
import Info from '@/lib/pages/home/components/Info';
import { Card } from '@chakra-ui/react/card';
import { Image } from '@chakra-ui/react/image';

const Home = () => {
  const { colorMode } = useColorMode();

  return (
    <>
      <Card.Root
        flexDirection={{ base: 'column', md: 'row' }}
        mb="4"
        variant="outline"
        className="page-wrapper"
      >
        <Card.Header style={{ paddingLeft: '2em', paddingTop: '.5em' }}>
          <Image
            style={{ height: '7em', width: '7em', alignSelf: 'center' }}
            src={colorMode === 'light' ? tridentLight : trident}
            alt="course logo"
          />
          <Card.Title>CIS Web Developer Certificate</Card.Title>
          <Card.Description className="outline">
            Program Requirements (44 credits)
          </Card.Description>
        </Card.Header>
        <Card.Body className="outline">
          <h3 className="outline">OUTCOMES</h3>
          <li>Build and maintain websites</li>
          <li>Work with stakeholders to create websites</li>
          <li>
            Research, assess, and appropriately apply emerging technology to
            support websites as needed in industry
          </li>
          <li>
            Comply with the ethics related to the use of copyrighted materials
            and intellectual property rights
          </li>
          <li>
            Demonstrate an entrepreneurial approach to web development sites and
            pages
          </li>
          <li>
            Manage career goals through creating effective resumes/CVs,
            developing interviewing skills, and setting goals
          </li>
        </Card.Body>
      </Card.Root>
      <SimpleGrid columns={{ base: 1, md: 2 }} gap={2}>
        {DATA.map((item) => {
          return <Info key={JSON.stringify(item)} item={item} />;
        })}
      </SimpleGrid>
    </>
  );
};

export default Home;
