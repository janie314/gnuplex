import {
  getKeyValue,
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
} from "@nextui-org/react";
import { useMemo } from "react";
import { APICall } from "../../lib/API.ts";

function TaggingTab(
  props: { closeHook: () => void; medialist: string[] },
) {
  const tablelist = useMemo(
    () =>
      props.medialist.map((mediafile: string, i: number) => (
        {
          key: i,
          name: mediafile.split("/").slice(-1).join(""),
          path: mediafile,
        }
      )),
    [props.medialist],
  );

  const columns = [
    {
      key: "name",
      label: "Name",
    },
  ];

  return (
    <>
      <Table
        color="primary"
        isStriped
        selectionMode="single"
        onRowAction={(key: any) => APICall.setMedia(tablelist[key].path)}
      >
        <TableHeader columns={columns}>
          {(column) => (
            <TableColumn key={column.key}>{column.label}</TableColumn>
          )}
        </TableHeader>
        <TableBody items={tablelist}>
          {(item) => (
            <TableRow
              key={item.key}
            >
              {(columnKey) => (
                <TableCell>
                  {getKeyValue(item, columnKey)}
                </TableCell>
              )}
            </TableRow>
          )}
        </TableBody>
      </Table>
    </>
  );
}

export { TaggingTab };
