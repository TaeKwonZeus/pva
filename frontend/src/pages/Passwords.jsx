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
        name,
      }),
    });

    if (!res.ok) {
      alert(res.statusText + " " + (await res.text()));
      return;
    }
    window.location.reload();
  }

  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <Button style={{ width: "200px" }}>
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
            <Button onClick={createVault} disabled={name === ""}>
              Create vault
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
      method: "PATCH",
      body: JSON.stringify({
        name,
        description,
        password,
      }),
    });

    if (!res.ok) {
      alert(res.statusText + " " + (await res.text()));
      return;
    }
    window.location.reload();
  }

  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <IconButton size="1" title="Edit this password" variant="soft">
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
                  size="1"
                  variant="ghost"
                  onClick={() => setVisible(!visible)}
                >
                  {!visible ? <EyeOpenIcon /> : <EyeClosedIcon />}
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
      alert(res.statusText + " " + (await res.text()));
      return;
    }

    window.location.reload();
  }

  return (
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <IconButton size="1" color="red" title="Delete this password">
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
      <Flex direction="column" ml="5" gap="2">
        <Flex align="center" gap="4" {...otherProps}>
          <Flex width="15px" height="15px" align="center" justify="center">
            <LockClosedIcon />
          </Flex>
          <Box width="200px">{password.name}</Box>
          <Box width="200px">{password.description}</Box>
          <Flex gap="2">
            <IconButton
              size="1"
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
      </Flex>
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
        name,
        description,
        password,
      }),
    });
    if (!res.ok) {
      if (res.status === 409)
        alert("Password with this name already exists in this vault!");
      else alert(res.statusText + " " + (await res.text()));
      return;
    }

    window.location.reload();
  }

  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <IconButton size="1" title="Add a new password">
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
                  size="1"
                  variant="ghost"
                  onClick={() => setVisible(!visible)}
                >
                  {!visible ? <EyeOpenIcon /> : <EyeClosedIcon />}
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
              disabled={name === "" || description === "" || password === ""}
            >
              Create password
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
      method: "PATCH",
      body: JSON.stringify({
        name,
      }),
    });
    setName("");
    if (!res.ok) {
      alert(res.statusText + " " + (await res.text()));
      return;
    }
    window.location.reload();
  }

  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <IconButton size="1" variant="soft" title="Edit this vault">
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
            <Button onClick={editVault} disabled={name === ""}>
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
      alert(res.statusText + " " + (await res.text()));
      return;
    }

    window.location.reload();
  }

  return (
    <AlertDialog.Root>
      <AlertDialog.Trigger>
        <IconButton size="1" color="red" title="Delete this vault">
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

function ShareVaultDialog({ vaultId }) {
  const [name, setName] = useState("");

  async function shareVault() {
    const res = await fetch(
      `/api/vaults/${vaultId}/share?target=${encodeURI(name)}`,
      {
        method: "POST",
      },
    );

    if (!res.ok) {
      alert(res.statusText + " " + (await res.text()));
      return;
    }

    alert("Vault shared successfully!");
  }

  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <IconButton size="1" title="Share this vault" variant="surface">
          <Share1Icon />
        </IconButton>
      </Dialog.Trigger>
      <Dialog.Content>
        <Dialog.Title>Share this vault</Dialog.Title>
        <Dialog.Description>
          Share this vault with another user.
        </Dialog.Description>
        <Flex direction="column" gap="3">
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Username
            </Text>
            <TextField.Root
              onChange={(e) => setName(e.target.value.trim())}
              placeholder="Username goes here"
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
            <Button onClick={shareVault} disabled={name === ""}>
              Share this vault
            </Button>
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
        <IconButton
          size="1"
          variant="ghost"
          onClick={() => setIsExpanded(!isExpanded)}
        >
          {isExpanded ? <ChevronUpIcon /> : <ChevronDownIcon />}
        </IconButton>
        <Box width="200px">{vault.name}</Box>
        <Box width="200px">{(vault?.passwords ?? []).length}</Box>
        <Flex gap="2">
          <CreatePasswordDialog vaultId={vault.id} />
          <ShareVaultDialog vaultId={vault.id} />
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
      .catch(alert);
  }, []);

  return (
    <Flex direction="column" gap="3">
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
