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
import {
  PlusIcon,
  ChevronDownIcon,
  Pencil2Icon,
  TrashIcon,
  ChevronUpIcon,
} from "@radix-ui/react-icons";
import { useState } from "react";

function ExpandableRow({ children, expandComponent, ...otherProps }) {
  const [isExpanded, setIsExpanded] = useState(false);
  return (
    <>
      <Table.Row {...otherProps} align="center">
        <Table.Cell>
          <IconButton
            variant="ghost"
            onClick={() => setIsExpanded(!isExpanded)}
          >
            {isExpanded ? <ChevronUpIcon /> : <ChevronDownIcon />}
          </IconButton>
        </Table.Cell>
        {children}
      </Table.Row>
      {isExpanded && (
        <Table.Row>
          <Table.Cell />
          {expandComponent}
        </Table.Row>
      )}
    </>
  );
}

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
            <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {Array(3).fill(
            // <Table.Row key="1" align="center">
            //   <Table.RowHeaderCell>
            //     <IconButton variant="ghost">
            //       <ChevronDownIcon />
            //     </IconButton>
            //   </Table.RowHeaderCell>
            //   <Table.Cell>SOME PASSWORDS</Table.Cell>
            //   <Table.Cell>ДОХУЯ</Table.Cell>
            //   <Table.Cell>
            //     <Flex gap="2">
            //       <IconButton title="Add a new password">
            //         <PlusIcon />
            //       </IconButton>
            //       <IconButton variant="soft" title="Edit this vault">
            //         <Pencil2Icon />
            //       </IconButton>
            //       <IconButton color="red" title="Delete this vault">
            //         <TrashIcon />
            //       </IconButton>
            //     </Flex>
            //   </Table.Cell>
            // </Table.Row>,
            <ExpandableRow
              key="1"
              expandComponent={<div>AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA</div>}
            >
              <Table.Cell>SOME PASSWORDS</Table.Cell>
              <Table.Cell>ДОХУЯ</Table.Cell>
              <Table.Cell>
                <Flex gap="2">
                  <IconButton title="Add a new password">
                    <PlusIcon />
                  </IconButton>
                  <IconButton variant="soft" title="Edit this vault">
                    <Pencil2Icon />
                  </IconButton>
                  <IconButton color="red" title="Delete this vault">
                    <TrashIcon />
                  </IconButton>
                </Flex>
              </Table.Cell>
            </ExpandableRow>,
          )}
        </Table.Body>
      </Table.Root>
    </Flex>
  );
}

export default Passwords;
