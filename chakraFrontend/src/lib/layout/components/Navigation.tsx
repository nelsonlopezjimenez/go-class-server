//@ts-nocheck
import { ColorModeButton } from '@/components/ui/color-mode';
import { LinkButton } from '@/components/ui/link-button';
import { SegmentedControl } from '@/components/ui/segmented-control';
import { Tooltip } from '@/components/ui/tooltip';
import { locationChecker } from '@/helpers/helpers';
import { version } from '@/version';
import { Flex } from '@chakra-ui/react/flex';
import { Float } from '@chakra-ui/react/float';
import { Icon } from '@chakra-ui/react/icon';
import { useEffect, useState } from 'react';
import { VscLinkExternal } from 'react-icons/vsc';
import { useLocation, useNavigate } from 'react-router';

const navLinks = {
  Home: '/',
  'Offline Links': '/offline-links',
  'Course Content': '/course-content',
  Calendar: '/calendar',
  'Lan Links': '/lan-links',
  About: '/about',
};

const Nav: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [youAreHere, setYouAreHere] = useState('Home');

  useEffect(() => {
    locationChecker(location, setYouAreHere);
  }, [location]);
  return (
    <>
      <Flex
        gap={4}
        // position="fixed"
        justify="space-evenly"
        alignItems="center"
        backdropFilter={'blur(5px)'}
        bgColor={{ _light: 'blue.900/50', _dark: 'blue.900/50' }}
        // top={0}
        width={'full'}
        px={4}
        height="14"
        boxSizing="border-box"
        as="nav"
        _light={{
          borderBottomWidth: '1px',
        }}
        shadow="md"
        zIndex={'docked'}
        borderBottomWidth="1px"
      >
        <Float offset="7">
          <ColorModeButton />
        </Float>
        {/* <Image
						id="logo"
						display={{ base: "none", md: "block" }}
						style={{ height: "2.5em" }}
						src={colorMode === "light" ? tridentLight : trident}
						alt="course logo"
            /> */}
        <Tooltip
          openDelay={10}
          closeDelay={100}
          content="You will be redirected to Gitea"
        >
          <LinkButton
            colorPalette="orange"
            variant="outline"
            display={{ base: 'none', md: 'flex' }}
            href="http://192.168.1.28:3000/CIS_Team_EDCC/MGLauncher/issues/new"
          >
            <Float placement="bottom-end">
              <Icon size="sm">
                <VscLinkExternal />
              </Icon>
            </Float>
            Submit a Bug Report
          </LinkButton>
        </Tooltip>
        <SegmentedControl
          colorPalette="blue"
          alignSelf="center"
          size={{ base: 'md', md: 'lg' }}
          defaultValue="Home"
          items={Object.keys(navLinks)}
          //@ts-ignore
          onValueChange={(e: Event) =>
            navigate(navLinks[e.value], { viewTransition: true })
          }
          value={youAreHere}
          // onMouseEnter={() => prefetchLinks(4, { force: true })}
        />
        <Flex
          id="version"
          display={{ base: 'none', md: 'block' }}
          // style={{ height: "2.5em" }}
        >
          v{version}
        </Flex>
      </Flex>
    </>
  );
};

export default Nav;
