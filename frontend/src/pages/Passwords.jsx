import { Button, Dialog, Flex, Text, TextField } from "@radix-ui/themes";
import { PlusIcon } from "@radix-ui/react-icons";

function Passwords() {
  return (
    <Dialog.Root>
      <Dialog.Trigger>
        <Button>
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
  );
}

export default Passwords;
