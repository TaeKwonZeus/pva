import {
  Box,
  Button,
  Dialog,
  Flex,
  Heading,
  IconButton,
  Separator,
  Text,
  TextField,
} from "@radix-ui/themes";
import {
  PlusIcon,
  ChevronDownIcon,
  TrashIcon,
  ChevronUpIcon,
  Pencil1Icon,
  LockClosedIcon,
  ClipboardCopyIcon,
} from "@radix-ui/react-icons";
import { useState } from "react";

function Vault({ vault }) {
  const [isExpanded, setIsExpanded] = useState(false);
  return (
    <>
      <Flex align="center" gap="4">
        <IconButton variant="ghost" onClick={() => setIsExpanded(!isExpanded)}>
          {isExpanded ? <ChevronUpIcon /> : <ChevronDownIcon />}
        </IconButton>
        <Box width="200px">{vault.name}</Box>
        <Box width="200px">{vault.passwords.length}</Box>
        <Flex gap="2">
          <IconButton>
            <PlusIcon />
          </IconButton>
          <IconButton variant="soft">
            <Pencil1Icon />
          </IconButton>
          <IconButton color="red">
            <TrashIcon />
          </IconButton>
        </Flex>
      </Flex>
      <Separator size="4" />
      {isExpanded &&
        vault.passwords.map((p) => (
          <>
            <Flex align="center" gap="4">
              <Flex width="15px" height="15px" align="center" justify="center">
                <LockClosedIcon />
              </Flex>
              <Box width="200px">{p.name}</Box>
              <Box width="250px">{p.description}</Box>
              <Flex gap="2">
                <IconButton variant="surface">
                  <ClipboardCopyIcon />
                </IconButton>
                <IconButton variant="soft">
                  <Pencil1Icon />
                </IconButton>
                <IconButton color="red">
                  <TrashIcon />
                </IconButton>
              </Flex>
            </Flex>
            <Separator size="4" />
          </>
        ))}
    </>
  );
}

function Passwords() {
  const vault = {
    name: "Example vault",
    passwords: [
      {
        name: "Example password",
        description: "Example description",
        password: "12345aboba",
        createdAt: "2012-07-04T18:10:00.000+09:00",
        updatedAt: "2012-07-04T18:10:00.000+09:00",
      },
      {
        name: "Example password",
        description: "Example description",
        password: "12345aboba",
        createdAt: "2012-07-04T18:10:00.000+09:00",
        updatedAt: "2012-07-04T18:10:00.000+09:00",
      },
    ],
  };

  return (
    <Flex direction="column" gap="3" width="800px">
      <Dialog.Root>
        <Dialog.Trigger>
          <Button style={{ width: "200px" }} mb="3">
            <PlusIcon /> New Vault
          </Button>
        </Dialog.Trigger>
        <Dialog.Content>
          <Dialog.Title>New Vault</Dialog.Title>
          <Dialog.Description>
            Create a new vault for secure password storage.
          </Dialog.Description>
          <Flex direction="column" gap="3">
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Name
              </Text>
              <TextField.Root placeholder="Vault name goes here" />
            </label>
          </Flex>
          <Flex gap="3" mt="4" justify="end">
            <Dialog.Close>
              <Button variant="soft" color="gray">
                Cancel
              </Button>
            </Dialog.Close>
            <Dialog.Close>
              <Button>Create</Button>
            </Dialog.Close>
          </Flex>
        </Dialog.Content>
      </Dialog.Root>

      <Heading>My Vaults</Heading>
      <Flex direction="column" gap="3">
        <Flex align="center" gap="4">
          <Box width="15px" />
          <Box width="200px">
            <Text weight="bold">Name</Text>
          </Box>
          <Box width="200px">
            <Text weight="bold">Password count</Text>
          </Box>
          <Box>
            <Text weight="bold">Actions</Text>
          </Box>
        </Flex>
        <Separator size="4"></Separator>
        {Array(3).fill(<Vault vault={vault} />)}
      </Flex>
    </Flex>
  );
}

export default Passwords;
