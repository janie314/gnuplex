import { useEffect, useState } from "react";
import { APICall } from "../../lib/API.ts";
import { WorkingSpinnerTSX } from "../WorkingSpinner.tsx";
import { Textarea } from "@nextui-org/react";

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
        <span className="subtitle">Media Directories</span>
        <Textarea
          value={mediadirs}
          onValueChange={setMediadirs}
          rows={10}
          placeholder="/mnt/externalssd/tv/twilight_zone/eye_of_the_beholder.av1"
        >
        </Textarea>
        <span className="subtitle">Excluded File Extensions</span>
        <Textarea
          value={fileExts}
          onValueChange={setFileExts}
          rows={10}
          placeholder=".pdf"
        >
        </Textarea>
        <div>
          <input
            type="button"
            value="Save Settings"
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
          />
          <input
            type="button"
            value="Refresh Library"
            onClick={() => {
              setRefreshLibraryWorking(true);
              APICall.setMediafiles().then(() =>
                setRefreshLibraryWorking(false)
              );
            }}
          />
          <WorkingSpinnerTSX
            visible={saveFileExtsWorking || saveMediadirsWorking ||
              refreshLibraryWorking}
          />
        </div>
        <div className="okcancel">
          <input
            type="button"
            value="OK"
            onClick={() => {
              props.closeHook();
            }}
          >
          </input>
        </div>
      </>
    );
  }
}

export { FilesAndDirectoriesTab };