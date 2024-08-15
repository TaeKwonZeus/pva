import { Theme, Container, Flex, Box, Separator } from "@radix-ui/themes";

function App(props) {
  return (
    <Theme appearance="dark">
      <Container>
        <header>
          <Flex align="center" justify="between" height="48px">
            <Box>Logo</Box>
            <Box>Navbar</Box>
            <Box>Details</Box>
          </Flex>
        </header>
        <main>{props.children}</main>
      </Container>
    </Theme>
  );
}

export default App;
