import { useMemo } from "react";
import { APICall } from "../lib/API.ts";
import "../App.css";
import {
  Card,
  CardBody,
  CardHeader,
  getKeyValue,
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
} from "@nextui-org/react";

function Medialist(
  props: {
    medialist: string[];
    subtitle: string;
  },
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
    <Card shadow="sm">
      <CardHeader>{props.subtitle}</CardHeader>
      <CardBody>
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
      </CardBody>
    </Card>
  );
}

export { Medialist };
