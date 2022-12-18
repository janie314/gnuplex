import { SyntheticEvent, useEffect, useState } from "react";
import "../App.css";
import { APICall } from "../lib/APICall";
import "./CRUDPopup.css";

function CRUDPopup(props: {
  visible: boolean;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const [mediadirs, setMediadirs] = useState("");
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
              const arr = mediadirs.trim().split("\n").filter((line) =>
                !/^\s*$/.test(line)
              ).map((line) => line.trim());
              APICall.setOriginMediadirs(arr);
            }}
          />
        </div>
        <div>
          <input
            type="button"
            value="Refresh Library"
            onClick={() => {
              APICall.refreshOriginMediafiles();
            }}
          />
        </div>
        <div className="okcancel">
          <input
            type="button"
            value="OK"
            onClick={() => props.setMediadirInputPopup(false)}
          >
          </input>
        </div>
      </div>
    );
  }
  return <></>;
}

export { CRUDPopup };
