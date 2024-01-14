import { useMemo } from "react";
import { APICall } from "../lib/API.ts";
import "../App.css";
import {
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
    setMedia: React.Dispatch<React.SetStateAction<string>>;
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
    <div className="moviegroup">
      <span className="subtitle">{props.subtitle}</span>
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
    </div>
  );
}

export { Medialist };
