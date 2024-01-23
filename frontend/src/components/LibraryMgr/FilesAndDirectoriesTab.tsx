import { useEffect, useState } from "react";
import { APICall } from "../../lib/API.ts";
import { WorkingSpinnerTSX } from "../WorkingSpinner.tsx";
import { Textarea } from "@nextui-org/react";
import { Button } from "@nextui-org/react";

function FilesAndDirectoriesTab(
  props: { selected: boolean; closeHook: () => void },
) {
  const [mediadirs, setMediadirs] = useState("");
  const [fileExts, setFileExts] = useState("");
  const [refreshLibraryWorking, setRefreshLibraryWorking] = useState(false);
  const [saveMediadirsWorking, setSaveMediadirsWorking] = useState(false);
  const [saveFileExtsWorking, setSaveFileExtsWorking] = useState(false);

  useEffect(() => {
    APICall.mediadirs().then((res) => {
      setMediadirs(res.join("\n"));
    });
    APICall.fileExts().then((res) => {
      setFileExts(res.join("\n"));
    });
  }, [props.selected]);

  if (props.selected) {
    return (
      <>
        <span className="text-lg font-bold">Media Directories</span>
        <Textarea
          value={mediadirs}
          onValueChange={setMediadirs}
          minRows={10}
          placeholder="/mnt/externalssd/tv/twilight_zone/eye_of_the_beholder.av1"
        >
        </Textarea>
        <span className="text-lg font-bold">Excluded File Extensions</span>
        <Textarea
          value={fileExts}
          onValueChange={setFileExts}
          minRows={10}
          placeholder=".pdf"
        >
        </Textarea>
        <div className="flex flex-row space-x-2">
          <Button
            size="sm"
            color="primary"
            onClick={() => {
              setSaveMediadirsWorking(true);
              setSaveFileExtsWorking(true);
              const arr1 = mediadirs.trim().split("\n").filter((
                line: string,
              ) => !/^\s*$/.test(line)).map((line: string) => line.trim());
              const arr2 = fileExts.trim().split("\n").filter((line: string) =>
                !/^\s*$/.test(line)
              ).map((line: string) => line.trim());
              APICall.setMediadirs(arr1).then(() =>
                setSaveMediadirsWorking(false)
              );
              APICall.setFileExts(arr2).then(() =>
                setSaveFileExtsWorking(false)
              );
            }}
          >
            Save Settings
          </Button>
          <Button
            size="sm"
            color="primary"
            onClick={() => {
              setRefreshLibraryWorking(true);
              APICall.setMediafiles().then(() =>
                setRefreshLibraryWorking(false)
              );
            }}
          >
            Refresh Library
          </Button>
          <WorkingSpinnerTSX
            visible={saveFileExtsWorking || saveMediadirsWorking ||
              refreshLibraryWorking}
          />
        </div>
        <div className="okcancel">
          <Button
            size="sm"
            color="primary"
            onClick={() => {
              props.closeHook();
            }}
          >
            OK
          </Button>
        </div>
      </>
    );
  }
}

export { FilesAndDirectoriesTab };
