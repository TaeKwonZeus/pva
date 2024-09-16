import { useEffect, useState } from "react";
import {
  AlertDialog,
  Box,
  Button,
  Code,
  Dialog,
  Flex,
  IconButton,
  ScrollArea,
  Table,
  Text,
  TextArea,
  TextField,
} from "@radix-ui/themes";
import greenCircle from "../assets/green-circle.png";
import grayCircle from "../assets/gray-circle.png";
import { Pencil1Icon, PlusIcon, TrashIcon } from "@radix-ui/react-icons";

function CreateDeviceDialog() {
  const [ip, setIp] = useState("");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  async function confirm() {
    const res = await fetch("/api/devices", {
      method: "POST",
      body: JSON.stringify({
        ip,
        name,
        description,
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
        <Button
          title="Add a device that wasn't discovered automatically"
          style={{ width: "200px" }}
        >
          <PlusIcon /> New Device
        </Button>
      </Dialog.Trigger>
      <Dialog.Content>
        <Dialog.Title>Update device</Dialog.Title>
        <Dialog.Description>Make changes to this device.</Dialog.Description>
        <Flex direction="column" gap="3">
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              IP
            </Text>
            <TextField.Root
              defaultValue={ip}
              onChange={(e) => setIp(e.target.value.trim())}
              placeholder="IP goes here"
            />
          </label>
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Name
            </Text>
            <TextField.Root
              defaultValue={name}
              onChange={(e) => setName(e.target.value.trim())}
              placeholder="Name goes here"
            />
          </label>
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Description
            </Text>
            <TextArea
              defaultValue={description}
              onChange={(e) => setDescription(e.target.value.trim())}
              placeholder="Description goes here"
            />
          </label>
          <Flex gap="3" justify="end">
            <Dialog.Close>
              <Button variant="soft" color="gray">
                Cancel
              </Button>
            </Dialog.Close>
            <Dialog.Close>
              <Button onClick={confirm} disabled={ip === "" || name === ""}>
                Create device
              </Button>
            </Dialog.Close>
          </Flex>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  );
}

function UpdateDeviceDialog({ device }) {
  const [ip, setIp] = useState(device.ip);
  const [name, setName] = useState(device.name);
  const [description, setDescription] = useState(device.description);

  async function confirm() {
    const res = await fetch(`/api/devices`, {
      method: "PUT",
      body: JSON.stringify({
        id: device.id ?? 0,
        ip,
        name,
        description,
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
        <IconButton title="Edit this device" size="1">
          <Pencil1Icon />
        </IconButton>
      </Dialog.Trigger>
      <Dialog.Content>
        <Dialog.Title>Update device</Dialog.Title>
        <Dialog.Description>Make changes to this device.</Dialog.Description>
        <Flex direction="column" gap="3">
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              IP
            </Text>
            <TextField.Root
              defaultValue={ip}
              onChange={(e) => setIp(e.target.value.trim())}
              placeholder="IP goes here"
            />
          </label>
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Name
            </Text>
            <TextField.Root
              defaultValue={name}
              onChange={(e) => setName(e.target.value.trim())}
              placeholder="Name goes here"
            />
          </label>
          <label>
            <Text as="div" size="2" mb="1" weight="bold">
              Description
            </Text>
            <TextArea
              defaultValue={description}
              onChange={(e) => setDescription(e.target.value.trim())}
              placeholder="Description goes here"
            />
          </label>
          <Flex gap="3" justify="end">
            <Dialog.Close>
              <Button variant="soft" color="gray">
                Cancel
              </Button>
            </Dialog.Close>
            <Dialog.Close>
              <Button
                onClick={confirm}
                disabled={
                  ip === device.ip &&
                  name === device.name &&
                  description === device.description
                }
              >
                Update device
              </Button>
            </Dialog.Close>
          </Flex>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  );
}

function DeleteDeviceDialog({ device }) {
  async function confirm() {
    const res = await fetch(`/api/devices?id=${device.id}`, {
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
        <IconButton title="Delete this device" color="red" size="1">
          <TrashIcon />
        </IconButton>
      </AlertDialog.Trigger>
      <AlertDialog.Content>
        <AlertDialog.Title>Delete device</AlertDialog.Title>
        <AlertDialog.Description>
          Are you sure? If this device is discoverable, it will still show up,
          but the name and description will be gone.
        </AlertDialog.Description>
        <Flex gap="3" mt="4" justify="end">
          <AlertDialog.Cancel>
            <Button variant="soft" color="gray">
              Cancel
            </Button>
          </AlertDialog.Cancel>
          <AlertDialog.Action onClick={confirm}>
            <Button color="red">Delete this device</Button>
          </AlertDialog.Action>
        </Flex>
      </AlertDialog.Content>
    </AlertDialog.Root>
  );
}

function Devices() {
  const [devices, setDevices] = useState([]);

  useEffect(() => {
    fetch("/api/devices")
      .then((r) => r.json())
      .then((v) => setDevices(v ?? []))
      .catch(alert);
  }, []);

  return (
    <Flex direction="column" gap="3">
      <CreateDeviceDialog />
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.ColumnHeaderCell width="100px">IP</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell width="100px">Name</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell width="250px">
              Description
            </Table.ColumnHeaderCell>
            {/*<Table.ColumnHeaderCell>Network Name</Table.ColumnHeaderCell>*/}
            {/*<Table.ColumnHeaderCell>MAC</Table.ColumnHeaderCell>*/}
            <Table.ColumnHeaderCell>Discoverable</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {devices.map((device, idx) => (
            <Table.Row key={idx}>
              <Table.Cell>
                <Code variant="ghost">{device.ip}</Code>
              </Table.Cell>
              <Table.Cell>
                {device.name !== "" ? (
                  device.name
                ) : (
                  <Text color="gray">None</Text>
                )}
              </Table.Cell>
              <Table.Cell>
                <ScrollArea
                  scrollbars="vertical"
                  type="auto"
                  style={{ maxHeight: "50px" }}
                >
                  <Box pr="1">
                    {device.description !== "" ? (
                      device.description
                    ) : (
                      <Text color="gray">None</Text>
                    )}
                  </Box>
                </ScrollArea>
              </Table.Cell>
              {/*<Table.Cell>{device.networkName}</Table.Cell>*/}
              {/*<Table.Cell>{device.mac}</Table.Cell>*/}
              <Table.Cell>
                <Flex align="center">
                  {device.connected ? (
                    <img width="21px" alt={"True"} src={greenCircle} />
                  ) : (
                    <img width="21px" alt={"True"} src={grayCircle} />
                  )}
                </Flex>
              </Table.Cell>
              <Table.Cell>
                <Flex gap="2">
                  <UpdateDeviceDialog device={device} />
                  <DeleteDeviceDialog device={device} />
                </Flex>
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table.Root>
    </Flex>
  );
}

export default Devices;
