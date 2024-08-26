import {
  Button,
  Dialog,
  Flex,
  Heading,
  IconButton,
  Table,
  Text,
  TextField,
} from "@radix-ui/themes";
import { PlusIcon, ChevronDownIcon } from "@radix-ui/react-icons";

function Passwords() {
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
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.ColumnHeaderCell />
            <Table.ColumnHeaderCell>Name</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell>Password count</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell />
          </Table.Row>
        </Table.Header>
        <Table.Body>
          <Table.Row>
            <Table.Cell>
              <IconButton variant="soft">
                <ChevronDownIcon />
              </IconButton>
            </Table.Cell>
            <Table.RowHeaderCell>SOME PASSWORDS</Table.RowHeaderCell>
            <Table.Cell>ДОХУЯ</Table.Cell>
            <Table.Cell>
              <IconButton>
                <PlusIcon />
              </IconButton>
            </Table.Cell>
          </Table.Row>

          <Table.Row>
            <Table.Cell>
              <IconButton variant="soft">
                <ChevronDownIcon />
              </IconButton>
            </Table.Cell>
            <Table.RowHeaderCell>SOME PASSWORDS</Table.RowHeaderCell>
            <Table.Cell>ДОХУЯ</Table.Cell>
            <Table.Cell>
              <IconButton>
                <PlusIcon />
              </IconButton>
            </Table.Cell>
          </Table.Row>

          <Table.Row>
            <Table.Cell>
              <IconButton variant="soft">
                <ChevronDownIcon />
              </IconButton>
            </Table.Cell>
            <Table.RowHeaderCell>SOME PASSWORDS</Table.RowHeaderCell>
            <Table.Cell>ДОХУЯ</Table.Cell>
            <Table.Cell>
              <IconButton>
                <PlusIcon />
              </IconButton>
            </Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table.Root>
    </Flex>
  );
}

export default Passwords;
