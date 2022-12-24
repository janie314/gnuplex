import { useEffect, useState } from "react";
import { APICall } from "../lib/APICall";
import "../App.css";
import "./CRUDPopup.css";
import { WorkingSpinnerTSX } from "./WorkingSpinner";

function CRUDPopup(props: {
  visible: boolean;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  closeHook: () => void;
}) {
  const [mediadirs, setMediadirs] = useState("");
  const [file_exts, setFileExts] = useState("");
  const [refreshLibraryWorking, setRefreshLibraryWorking] = useState(false);
  const [saveMediadirsWorking, setSaveMediadirsWorking] = useState(false);
  const [saveFileExtsWorking, setSaveFileExtsWorking] = useState(false);

  useEffect(() => {
    APICall.getOriginMediadirs().then((res: string[]) => {
      setMediadirs(res.join("\n"));
    });
  }, [props.visible]);

  if (props.visible) {
    return (
      <div className="crudpopup">
        <span className="subtitle">Media Directories</span>
        <textarea
          value={mediadirs}
          onChange={(e) => setMediadirs(e.target.value)}
          className="crudpopup-textarea"
          rows={10}
          placeholder="/mnt/externalssd/tv/twilight_zone/eye_of_the_beholder.av1"
        >
        </textarea>
                <span className="subtitle">Excluded File Extensions</span>
        <textarea
          value={file_exts}
          onChange={(e) => setFileExts(e.target.value)}
          className="crudpopup-textarea"
          rows={10}
          placeholder=".pdf\n.txt"
        >
        </textarea>
        <div>
          <input
            type="button"
            value="Save Settings"
            onClick={() => {
              setSaveMediadirsWorking(true);
              setSaveFileExtsWorking(true);
              const arr1 = mediadirs.trim().split("\n").filter((line) =>
                !/^\s*$/.test(line)
              ).map((line) => line.trim());
              const arr2 = mediadirs.trim().split("\n").filter((line) =>
                !/^\s*$/.test(line)
              ).map((line) => line.trim());
              APICall.setOriginMediadirs(arr1).then(() =>
                setSaveMediadirsWorking(false)
              );
              APICall.setOriginFileExts(arr2).then(() =>
                setSaveFileExtsWorking(false)
              );
            }}
          />
        </div>
        <div>
          <input
            type="button"
            value="Refresh Library"
            onClick={() => {
              setRefreshLibraryWorking(true);
              APICall.refreshOriginMediafiles().then(() =>
                setRefreshLibraryWorking(false)
              );
            }}
          />
          <WorkingSpinnerTSX visible={refreshLibraryWorking} />
        </div>
        <div className="okcancel">
          <input
            type="button"
            value="OK"
            onClick={() => {
              props.closeHook();
              props.setMediadirInputPopup(false);
            }}
          >
          </input>
        </div>
      </div>
    );
  }
  return <></>;
}

export { CRUDPopup };
