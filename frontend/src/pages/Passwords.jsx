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
  AlertDialog,
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
  Share1Icon,
} from "@radix-ui/react-icons";
import { useEffect, useState } from "react";

function CreateVaultDialog() {
  const [name, setName] = useState("");

  async function createVault() {
    const res = await fetch("/api/vaults/new", {
      method: "POST",
      body: JSON.stringify({
        name: name.trim(),
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
            <Button onClick={createVault} disabled={name.trim() === ""}>
              Create
            </Button>
          </Dialog.Close>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  );
}

function EditPasswordDialog({ vaultId, passwordId }) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [password, setPassword] = useState("");

  const [visible, setVisible] = useState(false);

  async function editPassword() {
    const res = await fetch(`/api/vaults/${vaultId}/${passwordId}`, {
      method: "PUT",
      body: JSON.stringify({
        name,
        description,
        password,
      }),
    });

    if (!res.ok) {
      console.error(res.statusText + " " + (await res.text()));
      return;
    }
    window.location.reload();
  }

  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <IconButton title="Edit this password" variant="soft">
          <Pencil1Icon />
        </IconButton>
      </Dialog.Trigger>
      <Dialog.Content>
        <Dialog.Title>Edit this password</Dialog.Title>
        <Dialog.Description>
          Edit this password. Omitted fields will not be changed.
        </Dialog.Description>
        <Flex direction="column" gap="3">
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Name
            </Text>
            <TextField.Root
              onChange={(e) => setName(e.target.value.trim())}
              placeholder="Password name goes here"
            />
          </label>
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Name
            </Text>
            <TextArea
              onChange={(e) => setDescription(e.target.value.trim())}
              placeholder="Description goes here"
            />
          </label>
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Password
            </Text>
            <TextField.Root
              type={visible ? "text" : "password"}
              onChange={(e) => setPassword(e.target.value.trim())}
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
            <Button
              onClick={editPassword}
              disabled={name === "" && description === "" && password === ""}
            >
              Edit this password
            </Button>
          </Dialog.Close>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  );
}

function DeletePasswordDialog({ vaultId, passwordId }) {
  async function deletePassword() {
    const res = await fetch(`/api/vaults/${vaultId}/${passwordId}`, {
      method: "DELETE",
    });

    if (!res.ok) {
      console.error(res.statusText + " " + (await res.text()));
      return;
    }

    window.location.reload();
  }

  return (
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <IconButton color="red" title="Delete this password">
          <TrashIcon />
        </IconButton>
      </AlertDialog.Trigger>
      <AlertDialog.Content maxWidth="450px">
        <AlertDialog.Title>Delete this vault</AlertDialog.Title>
        <AlertDialog.Description size="2">
          Are you sure? This password will no longer be accessible by anyone
          with access to it.
        </AlertDialog.Description>
        <Flex gap="3" mt="4" justify="end">
          <AlertDialog.Cancel>
            <Button variant="soft" color="gray">
              Cancel
            </Button>
          </AlertDialog.Cancel>
          <AlertDialog.Action>
            <Button color="red" onClick={deletePassword}>
              Delete this password
            </Button>
          </AlertDialog.Action>
        </Flex>
      </AlertDialog.Content>
    </AlertDialog.Root>
  );
}

function Password({ password, vaultId, ...otherProps }) {
  return (
    <>
      <Flex align="center" gap="4" {...otherProps}>
        <Flex width="15px" height="15px" align="center" justify="center" ml="5">
          <LockClosedIcon />
        </Flex>
        <Box width="200px">{password.name}</Box>
        <Box width="200px">{password.description}</Box>
        <Flex gap="2">
          <IconButton
            variant="surface"
            onClick={() => navigator.clipboard.writeText(password.password)}
            title="Copy to clipboard"
          >
            <ClipboardCopyIcon />
          </IconButton>

          <EditPasswordDialog vaultId={vaultId} passwordId={password.id} />

          <DeletePasswordDialog vaultId={vaultId} passwordId={password.id} />
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
    const res = await fetch(`/api/vaults/${vaultId}/new`, {
      method: "POST",
      body: JSON.stringify({
        name: name.trim(),
        description: description.trim(),
        password: password.trim(),
      }),
    });
    if (!res.ok) {
      if (res.status === 409)
        console.error("Password with this name already exists in this vault!");
      else console.error(res.statusText + " " + (await res.text()));
      return;
    }

    window.location.reload();
  }

  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <IconButton title="Add a new password">
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
              Password
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
            <Button
              onClick={createPassword}
              disabled={name.trim() === "" || password.trim() === ""}
            >
              Create
            </Button>
          </Dialog.Close>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  );
}

function EditVaultDialog({ vaultId }) {
  const [name, setName] = useState("");

  async function editVault() {
    const res = await fetch(`/api/vaults/${vaultId}`, {
      method: "PUT",
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
        <IconButton variant="soft" title="Edit this vault">
          <Pencil1Icon />
        </IconButton>
      </Dialog.Trigger>
      <Dialog.Content>
        <Dialog.Title>Edit this vault</Dialog.Title>
        <Dialog.Description>
          Edit this vault. Omitted fields will not be updated.
        </Dialog.Description>
        <Flex direction="column" gap="3">
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Name
            </Text>
            <TextField.Root
              onChange={(e) => setName(e.target.value.trim())}
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
            <Button onClick={editVault} disabled={name.trim() === ""}>
              Edit this vault
            </Button>
          </Dialog.Close>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  );
}

function DeleteVaultDialog({ vaultId }) {
  async function deleteVault() {
    const res = await fetch(`/api/vaults/${vaultId}`, {
      method: "DELETE",
    });

    if (!res.ok) {
      console.error(res.statusText + " " + (await res.text()));
      return;
    }

    window.location.reload();
  }

  return (
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <IconButton color="red" title="Delete this vault">
          <TrashIcon />
        </IconButton>
      </AlertDialog.Trigger>
      <AlertDialog.Content maxWidth="450px">
        <AlertDialog.Title>Delete this vault</AlertDialog.Title>
        <AlertDialog.Description size="2">
          Are you sure? This vault and ALL passwords belonging to it will no
          longer be accessible by anyone with access to them.
        </AlertDialog.Description>
        <Flex gap="3" mt="4" justify="end">
          <AlertDialog.Cancel>
            <Button variant="soft" color="gray">
              Cancel
            </Button>
          </AlertDialog.Cancel>
          <AlertDialog.Action>
            <Button color="red" onClick={deleteVault}>
              Delete this vault
            </Button>
          </AlertDialog.Action>
        </Flex>
      </AlertDialog.Content>
    </AlertDialog.Root>
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
          <IconButton variant="surface" title="Share this vault">
            <Share1Icon />
          </IconButton>

          <EditVaultDialog vaultId={vault.id} />

          <DeleteVaultDialog vaultId={vault.id} />
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
