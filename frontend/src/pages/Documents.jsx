import { Box, Code, Flex, ScrollArea, Table, Text } from "@radix-ui/themes";
import greenCircle from "../assets/green-circle.png";
import grayCircle from "../assets/gray-circle.png";

function CreateDocumentDialog() {}

function Documents() {
  return (
    <Flex direction="column" gap="3">
      <CreateDocumentDialog />
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

export default Documents;
