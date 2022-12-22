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
  const [refreshLibraryWorking, setRefreshLibraryWorking] = useState(false);
  const [saveMediadirsWorking, setSaveMediadirsWorking] = useState(false);

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
        <div>
          <input
            type="button"
            value="Save Media Directories"
            onClick={() => {
              setSaveMediadirsWorking(true);
              const arr = mediadirs.trim().split("\n").filter((line) =>
                !/^\s*$/.test(line)
              ).map((line) => line.trim());
              APICall.setOriginMediadirs(arr).then(() =>
                setSaveMediadirsWorking(false)
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
