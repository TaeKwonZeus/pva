import {
  Box,
  Button,
  Dialog,
  Flex,
  Heading,
  IconButton,
  Separator,
  Text,
  TextArea,
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
  EyeOpenIcon,
  EyeClosedIcon,
} from "@radix-ui/react-icons";
import { useEffect, useState } from "react";

function CreateVaultDialog() {
  const [name, setName] = useState("");

  async function createVault() {
    const res = await fetch("/api/vaults/new", {
      method: "POST",
      body: JSON.stringify({
        name,
      }),
    });
    setName("");
    if (!res.ok) {
      console.error(res.statusText + " " + (await res.text()));
      return;
    }
    window.location.reload();
  }

  return (
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
            <TextField.Root
              onChange={(e) => setName(e.target.value)}
              placeholder="Vault name goes here"
            />
          </label>
        </Flex>
        <Flex gap="3" mt="4" justify="end">
          <Dialog.Close>
            <Button variant="soft" color="gray">
              Cancel
            </Button>
          </Dialog.Close>
          <Dialog.Close>
            <Button onClick={createVault}>Create</Button>
          </Dialog.Close>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  );
}

function Password({ password, ...otherProps }) {
  return (
    <>
      <Flex align="center" gap="4" {...otherProps}>
        <Flex width="15px" height="15px" align="center" justify="center">
          <LockClosedIcon />
        </Flex>
        <Box width="200px">{password.name}</Box>
        <Box width="250px">{password.description}</Box>
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
  );
}

function CreatePasswordDialog({ vaultId }) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [password, setPassword] = useState("");

  const [visible, setVisible] = useState(false);

  async function createPassword() {
    console.log(`${name} ${description} ${password} ${vaultId}`);
  }

  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <IconButton>
          <PlusIcon />
        </IconButton>
      </Dialog.Trigger>
      <Dialog.Content>
        <Dialog.Title>New Password</Dialog.Title>
        <Dialog.Description>Add a new password</Dialog.Description>
        <Flex direction="column" gap="3">
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Name
            </Text>
            <TextField.Root
              onChange={(e) => setName(e.target.value)}
              placeholder="Password name goes here"
            />
          </label>
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Name
            </Text>
            <TextArea
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Description goes here"
            />
          </label>
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Name
            </Text>
            <TextField.Root
              type={visible ? "text" : "password"}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Password goes here"
            >
              <TextField.Slot>
                <IconButton
                  variant="ghost"
                  onClick={() => setVisible(!visible)}
                >
                  {visible ? <EyeOpenIcon /> : <EyeClosedIcon />}
                </IconButton>
              </TextField.Slot>
            </TextField.Root>
          </label>
        </Flex>
        <Flex gap="3" mt="4" justify="end">
          <Dialog.Close>
            <Button variant="soft" color="gray">
              Cancel
            </Button>
          </Dialog.Close>
          <Dialog.Close>
            <Button onClick={createPassword}>Create</Button>
          </Dialog.Close>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  );
}

function Vault({ vault, ...otherProps }) {
  const [isExpanded, setIsExpanded] = useState(false);
  return (
    <>
      <Flex align="center" gap="4" {...otherProps}>
        <IconButton variant="ghost" onClick={() => setIsExpanded(!isExpanded)}>
          {isExpanded ? <ChevronUpIcon /> : <ChevronDownIcon />}
        </IconButton>
        <Box width="200px">{vault.name}</Box>
        <Box width="200px">{(vault?.passwords ?? []).length}</Box>
        <Flex gap="2">
          <CreatePasswordDialog vaultId={vault.id} />
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
        (vault?.passwords ?? []).map((p, index) => (
          <Password password={p} vaultId={vault.id} key={index} />
        ))}
    </>
  );
}

function Passwords() {
  const [vaults, setVaults] = useState([]);

  useEffect(() => {
    fetch("/api/vaults")
      .then((r) => r.json())
      .then((v) => setVaults(v))
      .catch(console.error);
  }, []);

  return (
    <Flex direction="column" gap="3" width="800px">
      <CreateVaultDialog />

      {vaults.length > 0 ? (
        <>
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
            {vaults.map((vault, index) => (
              <Vault vault={vault} key={index} />
            ))}
          </Flex>
        </>
      ) : (
        <Heading>Nothing to see here</Heading>
      )}
    </Flex>
  );
}

export default Passwords;
