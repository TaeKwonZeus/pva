import {
  Box,
  Button,
  Container,
  Flex,
  Heading,
  Link,
  TextField,
} from "@radix-ui/themes";
import {
  FileIcon,
  HomeIcon,
  LockClosedIcon,
  MagnifyingGlassIcon,
  SketchLogoIcon,
} from "@radix-ui/react-icons";
import { Outlet, useNavigate } from "react-router-dom";
import { logOut } from "./auth.js";

function App() {
  const navigate = useNavigate();

  const links = [
    {
      href: "#",
      label: "Home",
      icon: <HomeIcon />,
    },
    {
      href: "/passwords",
      label: "Passwords",
      icon: <LockClosedIcon />,
    },
    {
      href: "#",
      label: "Documents",
      icon: <FileIcon />,
    },
  ];

  return (
    <Container mx="3">
      <Flex align="center" justify="between" height="48px" mb="3">
        <Flex align="center" justify="center" gap="1">
          <SketchLogoIcon width="24px" height="24px" />
          <Heading>PVA</Heading>
        </Flex>
        <Box>
          <TextField.Root
            radius="full"
            variant="surface"
            style={{ width: 400 }}
            tabIndex={1}
            placeholder="Search"
          >
            <TextField.Slot>
              <MagnifyingGlassIcon />
            </TextField.Slot>
          </TextField.Root>
        </Box>
        <Button
          variant="ghost"
          color="gray"
          highContrast
          onClick={async () => {
            await logOut();
            navigate("/auth");
          }}
        >
          Log Out
        </Button>
      </Flex>
      <Flex>
        <Flex direction="column" gap="3" width="150px" mr="3">
          {links.map(({ href, label, icon }, idx) => (
            <Link
              key={idx}
              href={href}
              underline="hover"
              color="gray"
              weight="bold"
              highContrast
              size="4"
            >
              <Flex align="center" gap="1">
                {icon}
                {label}
              </Flex>
            </Link>
          ))}
        </Flex>
        <Outlet />
      </Flex>
    </Container>
  );
}

export default App;
