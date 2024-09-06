import {
  Box,
  Button,
  Container,
  Flex,
  Heading,
  HoverCard,
  Link,
  Separator,
  Text,
  TextField,
} from "@radix-ui/themes";
import {
  ArchiveIcon,
  FileIcon,
  HomeIcon,
  LockClosedIcon,
  MagnifyingGlassIcon,
  SketchLogoIcon,
} from "@radix-ui/react-icons";
import { Outlet, useNavigate } from "react-router-dom";
import { logOut } from "./auth.js";
import { useState } from "react";

const searchEntryIcons = {
  vault: <ArchiveIcon />,
  password: <LockClosedIcon />,
};

function SearchEntries({ index, text }) {
  const filtered = index.filter((item) =>
    item.title.toLowerCase().trim().includes(text.toLowerCase().trim()),
  );
  if (filtered.length === 0) {
    return <Text size="2">{"Didn't find anything :P"}</Text>;
  }

  return filtered.map((entry, i) => (
    <>
      <Link href={entry.url}>
        <Flex gap="2" key={i} align="center">
          <Box>{searchEntryIcons[entry["type"]]}</Box>
          <Text size="2">{entry.title}</Text>
        </Flex>
      </Link>

      {i !== filtered.length - 1 && <Separator size="4" />}
    </>
  ));
}

function SearchBar() {
  const [index, setIndex] = useState([]);
  const [loaded, setLoaded] = useState(false);

  async function loadIndex() {
    if (loaded) return;

    const res = await fetch("/api/index");
    if (!res.ok) {
      alert(res.statusText + " " + (await res.text()));
      return;
    }

    setIndex(await res.json());
    setLoaded(true);
  }

  const [focused, setFocused] = useState(false);
  const [text, setText] = useState("");

  return (
    <Flex onMouseEnter={loadIndex} direction="column">
      <TextField.Root
        radius="full"
        variant="surface"
        style={{ width: "400px" }}
        tabIndex={1}
        placeholder="Search"
        onChange={(e) => setText(e.target.value)}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
      >
        <TextField.Slot>
          <MagnifyingGlassIcon />
        </TextField.Slot>
      </TextField.Root>
      {/* TODO fix positioning */}
      <HoverCard.Root open={focused && text !== ""}>
        <HoverCard.Trigger>
          <Button style={{ width: 0, height: 0 }} />
        </HoverCard.Trigger>
        <HoverCard.Content width="400px">
          <Flex direction="column" gap="2">
            <SearchEntries index={index} text={text} />
          </Flex>
        </HoverCard.Content>
      </HoverCard.Root>
    </Flex>
  );
}

function App() {
  const navigate = useNavigate();

  const links = [
    {
      href: "/",
      label: "Home",
      icon: <HomeIcon />,
    },
    {
      href: "/passwords",
      label: "Passwords",
      icon: <LockClosedIcon />,
    },
    {
      href: "/",
      label: "Documents",
      icon: <FileIcon />,
    },
  ];

  return (
    <Container mx="3">
      <Flex align="center" justify="between" height="48px" mb="3">
        <Link href="/" color="gray" highContrast underline="none">
          <Flex align="center" justify="center" gap="1">
            <SketchLogoIcon width="24px" height="24px" />
            <Heading>PVA</Heading>
          </Flex>
        </Link>
        <SearchBar />
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
